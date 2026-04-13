import { useState, useEffect, useCallback, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { getAuth } from '../store/authStore'
import {
  getGame,
  makeMove,
  resignGame,
  offerDraw,
  acceptDraw,
  declineDraw,
  type GameResponse,
  type DrawOfferResponse,
} from '../api/chessApi'
import { useWebSocket } from '../hooks/useWebSocket'
import { ChessBoard } from '../components/ChessBoard'
import { PlayerProfile } from '../components/PlayerProfile'
import {
  GameStatus,
  GameResult,
  GameEndReason,
  RESULT_LABELS,
  END_REASON_LABELS,
} from '../lib/game'
import { Color } from '../lib/chess'
import { SERVER_EVENT, CLIENT_TYPE, KEY, type ServerMessage } from '../lib/wsMessages'
import * as S from './GameView.styled'
import { PrimaryButton, SecondaryButton, DangerButton } from '../styles/styled'

const STATUS_LABEL_YOUR_TURN = 'Your turn'
const STATUS_LABEL_OPPONENT_TURN = "Opponent's turn"
const STATUS_LABEL_WAITING = 'Waiting for player'
const PROMOTION_SUFFIX = '=Q'
const DRAW_OFFER_TIMEOUT_SEC = 30

export function GameView() {
  const { id } = useParams<{ id: string }>()!
  const [game, setGame] = useState<GameResponse | null>(null)
  const [errorMsg, setErrorMsg] = useState('')
  const [loading, setLoading] = useState(true)
  const [presentPlayers, setPresentPlayers] = useState<Set<string>>(new Set())
  const [drawOfferFromOpponent, setDrawOfferFromOpponent] = useState<DrawOfferResponse | null>(null)
  const [hasOfferedDraw, setHasOfferedDraw] = useState(false)
  const [drawOfferCountdown, setDrawOfferCountdown] = useState<number | null>(null)
  const [clockTick, setClockTick] = useState(0)
  const drawOfferTimerRef = useRef<ReturnType<typeof setInterval> | null>(null)
  const clockTimerRef = useRef<ReturnType<typeof setInterval> | null>(null)
  const timeoutTriggered = useRef(false)
  const navigate = useNavigate()

  useEffect(() => {
    if (!game?.time_control || game.time_control.base === 0 || game.status !== 'active') return
    clockTimerRef.current = setInterval(() => setClockTick(t => t + 1), 1000)
    return () => {
      if (clockTimerRef.current) clearInterval(clockTimerRef.current)
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [game?.time_control?.base, game?.status])

  const loadGame = useCallback(async () => {
    if (!id) return
    try {
      setGame(await getGame(id))
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : 'Failed to load game')
    } finally {
      setLoading(false)
    }
  }, [id])

  useEffect(() => {
    void loadGame()
  }, [loadGame])

  const stopDrawOfferCountdown = useCallback(() => {
    if (drawOfferTimerRef.current) {
      clearInterval(drawOfferTimerRef.current)
      drawOfferTimerRef.current = null
    }
    setDrawOfferCountdown(null)
  }, [])

  const startDrawOfferCountdown = useCallback(() => {
    stopDrawOfferCountdown()
    setDrawOfferCountdown(DRAW_OFFER_TIMEOUT_SEC)
    drawOfferTimerRef.current = setInterval(() => {
      setDrawOfferCountdown(prev => {
        if (prev === null || prev <= 1) {
          stopDrawOfferCountdown()
          setDrawOfferFromOpponent(null)
          return null
        }
        return prev - 1
      })
    }, 1000)
  }, [stopDrawOfferCountdown])

  useEffect(() => {
    return () => {
      stopDrawOfferCountdown()
    }
  }, [stopDrawOfferCountdown])

  const handleWsMessage = useCallback(
    (data: unknown) => {
      const msg = data as ServerMessage
      const msgType = msg[KEY.TYPE] as string | undefined

      if (msgType === SERVER_EVENT.PRESENCE && msg[KEY.GAME_ID] === id) {
        setPresentPlayers(new Set((msg[KEY.PLAYERS] as string[]) ?? []))
        return
      }
      if (
        (msgType === SERVER_EVENT.GAME_JOINED ||
          msgType === SERVER_EVENT.MOVE_MADE ||
          msgType === SERVER_EVENT.GAME_RESIGNED ||
          msgType === SERVER_EVENT.DRAW_ACCEPTED ||
          msgType === SERVER_EVENT.GAME_TIMED_OUT) &&
        msg[KEY.GAME_ID] === id
      ) {
        if (
          msgType === SERVER_EVENT.GAME_RESIGNED ||
          msgType === SERVER_EVENT.DRAW_ACCEPTED ||
          msgType === SERVER_EVENT.GAME_TIMED_OUT
        ) {
          setDrawOfferFromOpponent(null)
          setHasOfferedDraw(false)
          stopDrawOfferCountdown()
        }
        void loadGame()
        return
      }
      if (msgType === SERVER_EVENT.DRAW_OFFERED && msg[KEY.GAME_ID] === id) {
        const currentUserId = getAuth().user?.id
        const offeredBy = msg[KEY.OFFERED_BY] as string | undefined
        if (offeredBy && offeredBy !== currentUserId) {
          setDrawOfferFromOpponent({ offered_by: offeredBy, offered_at: new Date().toISOString() })
          startDrawOfferCountdown()
        }
        return
      }
      if (
        (msgType === SERVER_EVENT.DRAW_DECLINED || msgType === SERVER_EVENT.DRAW_OFFER_EXPIRED) &&
        msg[KEY.GAME_ID] === id
      ) {
        setDrawOfferFromOpponent(null)
        setHasOfferedDraw(false)
        stopDrawOfferCountdown()
        return
      }
      if (msg[KEY.GAME_ID] === id) void loadGame()
    },
    [id, loadGame, startDrawOfferCountdown, stopDrawOfferCountdown],
  )

  const { connected, send } = useWebSocket(handleWsMessage)

  useEffect(() => {
    if (!id || !connected) return
    send({ type: CLIENT_TYPE.JOIN_GAME, data: { [KEY.GAME_ID]: id } })
    return () => {
      send({ type: CLIENT_TYPE.LEAVE_GAME, data: { [KEY.GAME_ID]: id } })
    }
  }, [id, connected, send])

  const handleMove = async (from: string, to: string) => {
    if (!id) return
    setErrorMsg('')
    try {
      if (hasOfferedDraw) {
        setHasOfferedDraw(false)
      }
      setDrawOfferFromOpponent(null)
      stopDrawOfferCountdown()

      const g = await makeMove(id, { from, to })
      setGame(g)
    } catch (err) {
      await loadGame()
      setErrorMsg(err instanceof Error ? err.message : 'Invalid move')
    }
  }

  const handleResign = async () => {
    if (!id) return
    if (!window.confirm('Are you sure you want to resign?')) return
    setErrorMsg('')
    try {
      const g = await resignGame(id)
      setGame(g)
      setHasOfferedDraw(false)
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : 'Failed to resign')
    }
  }

  const handleOfferDraw = async () => {
    if (!id) return
    setErrorMsg('')
    try {
      const g = await offerDraw(id)
      setGame(g)
      setHasOfferedDraw(true)
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : 'Failed to offer draw')
    }
  }

  const handleAcceptDraw = async () => {
    if (!id) return
    setErrorMsg('')
    try {
      const g = await acceptDraw(id)
      setGame(g)
      setDrawOfferFromOpponent(null)
      setHasOfferedDraw(false)
      stopDrawOfferCountdown()
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : 'Failed to accept draw')
    }
  }

  const handleDeclineDraw = async () => {
    if (!id) return
    setErrorMsg('')
    try {
      const g = await declineDraw(id)
      setGame(g)
      setDrawOfferFromOpponent(null)
      setHasOfferedDraw(false)
      stopDrawOfferCountdown()
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : 'Failed to decline draw')
    }
  }

  useEffect(() => {
    if (!game || game.status !== 'active') {
      timeoutTriggered.current = false
      return
    }
    if (!game.time_control || game.time_control.base === 0) return
    if (timeoutTriggered.current) return
    const hasTimedOut = (color: Color) => {
      if (!game.time_control || game.time_control.base === 0) return false
      const stored = color === Color.White ? game.white_remaining : game.black_remaining
      if (stored == null) return false
      if (!game.last_move_at) return false
      const elapsed = Date.now() - new Date(game.last_move_at).getTime()
      return stored - elapsed <= 0
    }
    if (hasTimedOut(Color.White) || hasTimedOut(Color.Black)) {
      timeoutTriggered.current = true
      void loadGame()
    }
  })

  if (loading)
    return (
      <S.GameViewWrapper>
        <p>Loading game...</p>
      </S.GameViewWrapper>
    )
  if (!game)
    return (
      <S.GameViewWrapper>
        <p>
          Game not found.{' '}
          <PrimaryButton onClick={() => navigate('/lobby')}>Back to Lobby</PrimaryButton>
        </p>
      </S.GameViewWrapper>
    )

  const currentUserId = getAuth().user?.id
  const isBlack = game.black_player?.id === currentUserId
  const isYourTurn =
    (game.turn === Color.White && game.white_player.id === currentUserId) ||
    (game.turn === Color.Black && game.black_player?.id === currentUserId)
  const isGameActive = game.status === GameStatus.Active
  const isPlayerInGame =
    game.white_player.id === currentUserId || game.black_player?.id === currentUserId

  const getLiveRemaining = (color: Color): number | null => {
    if (!game.time_control || game.time_control.base === 0) return null
    const stored = color === Color.White ? game.white_remaining : game.black_remaining
    if (stored == null) return null
    if (!isGameActive) return stored
    if (!game.last_move_at) return stored
    void clockTick
    if (color === game.turn) {
      const elapsed = Date.now() - new Date(game.last_move_at).getTime()
      const live = stored - elapsed
      return Math.max(0, live)
    }
    return stored
  }

  const formatClock = (ms: number | null): string => {
    if (ms == null) return ''
    const totalSec = Math.ceil(ms / 1000)
    const min = Math.floor(totalSec / 60)
    const sec = totalSec % 60
    if (min >= 60) {
      const h = Math.floor(min / 60)
      return `${h}:${String(min % 60).padStart(2, '0')}:${String(sec).padStart(2, '0')}`
    }
    return `${min}:${String(sec).padStart(2, '0')}`
  }

  const statusLabel = (): string => {
    if (game.status === GameStatus.Active) {
      return isYourTurn ? STATUS_LABEL_YOUR_TURN : STATUS_LABEL_OPPONENT_TURN
    }
    if (game.status === GameStatus.Finished && game.result) {
      return formatResult(game.result, game.end_reason)
    }
    return STATUS_LABEL_WAITING
  }

  const formatResult = (result: GameResult, endReason: GameEndReason | null): string => {
    const label = RESULT_LABELS[result]
    if (!endReason) return label
    if (endReason === GameEndReason.Timeout) {
      const loserColor = result === GameResult.WhiteWins ? 'Black' : 'White'
      return `${label} — ${loserColor} ran out of time`
    }
    return `${label} ${END_REASON_LABELS[endReason]}`
  }

  const topColor: Color = isBlack ? Color.White : Color.Black
  const bottomColor: Color = isBlack ? Color.Black : Color.White
  const topPlayer = isBlack ? game.white_player : game.black_player
  const bottomPlayer = isBlack ? game.black_player : game.white_player

  const topIsPresent = topPlayer ? presentPlayers.has(topPlayer.id) : false
  const bottomIsPresent = bottomPlayer ? presentPlayers.has(bottomPlayer.id) : false
  const topIsActive = game.status === GameStatus.Active && game.turn === topColor
  const bottomIsActive = game.status === GameStatus.Active && game.turn === bottomColor

  const moves = game.moves ?? []
  const movePairs = Array.from({ length: Math.ceil(moves.length / 2) }, (_, i) => ({
    num: i + 1,
    white: moves[i * 2],
    black: moves[i * 2 + 1],
  }))

  const formatMove = (from: string, to: string, promotion: string | null) =>
    `${from}${to}${promotion === 'queen' ? PROMOTION_SUFFIX : ''}`

  const userHasPendingOffer = hasOfferedDraw
  const opponentHasOffered = drawOfferFromOpponent !== null

  return (
    <S.GameViewWrapper>
      <S.GameContent>
        <S.BoardColumn>
          <S.GameHeaderBar>
            <SecondaryButton onClick={() => navigate('/lobby')}>← Lobby</SecondaryButton>
            <S.GameStatus $status={game.status}>{statusLabel()}</S.GameStatus>
            <S.WsIndicator title={`WebSocket ${connected ? 'connected' : 'disconnected'}`}>
              {connected ? '●' : '○'}
            </S.WsIndicator>
          </S.GameHeaderBar>

          {errorMsg && <S.ErrorBanner>{errorMsg}</S.ErrorBanner>}

          <S.BoardContainer>
            {topPlayer && (
              <PlayerProfile
                key={`top-${topPlayer.id}`}
                username={topPlayer.username}
                color={topColor}
                isActive={topIsActive}
                isPresent={topIsPresent}
                clock={formatClock(getLiveRemaining(topColor))}
              />
            )}

            <ChessBoard
              fen={game.fen}
              turn={game.turn}
              status={game.status}
              onMove={handleMove}
              flipped={isBlack}
            />

            {bottomPlayer && (
              <PlayerProfile
                key={`bottom-${bottomPlayer.id}`}
                username={bottomPlayer.username}
                color={bottomColor}
                isActive={bottomIsActive}
                isPresent={bottomIsPresent}
                clock={formatClock(getLiveRemaining(bottomColor))}
              />
            )}
          </S.BoardContainer>
        </S.BoardColumn>

        <S.GameSidebar>
          <S.MoveHistory>
            <h3>Moves ({moves.length})</h3>
            {moves.length === 0 ? (
              <S.NoMoves>No moves yet</S.NoMoves>
            ) : (
              <S.MovesList>
                {movePairs.map(({ num, white, black }) => (
                  <S.MoveRow key={num}>
                    <span className="move-num">{num}.</span>
                    <span className="move">
                      {white ? formatMove(white.from, white.to, white.promotion) : ''}
                    </span>
                    <span className="move">
                      {black ? formatMove(black.from, black.to, black.promotion) : ''}
                    </span>
                  </S.MoveRow>
                ))}
              </S.MovesList>
            )}
          </S.MoveHistory>

          {game.status === GameStatus.Waiting && !game.black_player && (
            <S.WaitingNotice>
              <p>Waiting for opponent to join...</p>
            </S.WaitingNotice>
          )}

          {game.status === GameStatus.Finished && (
            <S.GameOverNotice>
              <p>{statusLabel()}</p>
            </S.GameOverNotice>
          )}

          {opponentHasOffered && isGameActive && (
            <S.DrawOfferInline>
              <p className="draw-offer-text">
                Opponent offers a draw
                {drawOfferCountdown !== null && (
                  <span className="draw-offer-timer"> ({drawOfferCountdown}s)</span>
                )}
              </p>
              <S.DrawOfferButtons>
                <PrimaryButton onClick={handleAcceptDraw}>Accept</PrimaryButton>
                <SecondaryButton onClick={handleDeclineDraw}>Decline</SecondaryButton>
              </S.DrawOfferButtons>
            </S.DrawOfferInline>
          )}

          {isGameActive && isPlayerInGame && (
            <S.GameActions>
              {!userHasPendingOffer && (
                <SecondaryButton onClick={handleOfferDraw}>Offer Draw</SecondaryButton>
              )}
              {userHasPendingOffer && <S.DrawPendingLabel>Draw offer pending...</S.DrawPendingLabel>}
              <DangerButton onClick={handleResign}>Resign</DangerButton>
            </S.GameActions>
          )}
        </S.GameSidebar>
      </S.GameContent>
    </S.GameViewWrapper>
  )
}
