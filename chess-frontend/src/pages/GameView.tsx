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
      // Making a move implicitly declines any pending draw offer
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

  // Auto-detect timeout — when live clock hits 0, call server to end the game
  useEffect(() => {
    if (!game || game.status !== 'active') {
      timeoutTriggered.current = false
      return
    }
    if (!game.time_control || game.time_control.base === 0) return
    if (timeoutTriggered.current) return
    // Inline the clock check to avoid reference issues
    const hasTimedOut = (color: Color) => {
      if (!game.time_control || game.time_control.base === 0) return false
      const stored = color === Color.White ? game.white_remaining : game.black_remaining
      if (stored == null) return false
      if (color === game.turn && game.last_move_at) {
        const elapsed = Date.now() - new Date(game.last_move_at).getTime()
        return stored - elapsed <= 0
      }
      return false
    }
    if (hasTimedOut(Color.White) || hasTimedOut(Color.Black)) {
      timeoutTriggered.current = true
      void loadGame()
    }
  })

  if (loading)
    return (
      <div className="lobby">
        <p>Loading game...</p>
      </div>
    )
  if (!game)
    return (
      <div className="lobby">
        <p>
          Game not found.{' '}
          <button className="btn-primary" onClick={() => navigate('/lobby')}>
            Back to Lobby
          </button>
        </p>
      </div>
    )

  const currentUserId = getAuth().user?.id
  const isBlack = game.black_player?.id === currentUserId
  const isYourTurn =
    (game.turn === Color.White && game.white_player.id === currentUserId) ||
    (game.turn === Color.Black && game.black_player?.id === currentUserId)
  const isGameActive = game.status === GameStatus.Active
  const isPlayerInGame =
    game.white_player.id === currentUserId || game.black_player?.id === currentUserId

  // ─── Live clock calculation ────────────────────────────────────────────────
  const getLiveRemaining = (color: Color): number | null => {
    if (!game.time_control || game.time_control.base === 0) return null
    const stored = color === Color.White ? game.white_remaining : game.black_remaining
    if (stored == null) return null
    if (!isGameActive) return stored
    // Clock hasn't started yet — return full time
    if (!game.last_move_at) return stored
    // Deduct elapsed time for the player whose turn it is
    // clockTick forces a re-render every second so the live clock updates
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
      // Show who ran out of time explicitly
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
    <div className="game-view">
      <div className="game-content">
        <div className="board-column">
          <div className="game-header-bar">
            <button className="btn-secondary" onClick={() => navigate('/lobby')}>
              ← Lobby
            </button>
            <span className={`game-status ${game.status}`}>{statusLabel()}</span>
            <span
              className="ws-indicator"
              title={`WebSocket ${connected ? 'connected' : 'disconnected'}`}
            >
              {connected ? '●' : '○'}
            </span>
          </div>

          {errorMsg && <div className="error-banner">{errorMsg}</div>}

          <div className="board-container">
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
          </div>
        </div>

        <div className="game-sidebar">
          <div className="move-history">
            <h3>Moves ({moves.length})</h3>
            {moves.length === 0 ? (
              <p className="no-moves">No moves yet</p>
            ) : (
              <div className="moves-list">
                {movePairs.map(({ num, white, black }) => (
                  <div key={num} className="move-row">
                    <span className="move-num">{num}.</span>
                    <span className="move">
                      {white ? formatMove(white.from, white.to, white.promotion) : ''}
                    </span>
                    <span className="move">
                      {black ? formatMove(black.from, black.to, black.promotion) : ''}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </div>

          {game.status === GameStatus.Waiting && !game.black_player && (
            <div className="waiting-notice">
              <p>Waiting for opponent to join...</p>
              <p className="hint">Share the game link or wait from another account.</p>
            </div>
          )}

          {game.status === GameStatus.Finished && (
            <div className="game-over-notice">
              <p>{statusLabel()}</p>
            </div>
          )}

          {/* Inline draw offer from opponent (chess.com style — no modal) */}
          {opponentHasOffered && isGameActive && (
            <div className="draw-offer-inline">
              <p className="draw-offer-text">
                Opponent offers a draw
                {drawOfferCountdown !== null && (
                  <span className="draw-offer-timer"> ({drawOfferCountdown}s)</span>
                )}
              </p>
              <div className="draw-offer-buttons">
                <button className="btn-primary" onClick={handleAcceptDraw}>
                  Accept
                </button>
                <button className="btn-secondary" onClick={handleDeclineDraw}>
                  Decline
                </button>
              </div>
            </div>
          )}

          {/* Resign / Draw buttons */}
          {isGameActive && isPlayerInGame && (
            <div className="game-actions">
              {!userHasPendingOffer && (
                <button className="btn-secondary" onClick={handleOfferDraw}>
                  Offer Draw
                </button>
              )}
              {userHasPendingOffer && (
                <span className="draw-pending-label">Draw offer pending...</span>
              )}
              <button className="btn-danger" onClick={handleResign}>
                Resign
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
