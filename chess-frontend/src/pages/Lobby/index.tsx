import { useState, useEffect, useCallback, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { getAuth, clearAuth } from '../../store/authStore'
import {
  getWaitingGames,
  getMyGames,
  createGame,
  joinGame,
  type GameResponse,
  type CreateGameRequest,
} from '../../api/chessApi'
import { GameStatus } from '../../lib/game'
import { ALL_CATEGORIES } from '../../lib/game/timeControls'
import * as S from './styles'
import { PrimaryButton, SecondaryButton } from '../../styles/styled'

const STATUS_COLORS: Record<GameStatus, string> = {
  [GameStatus.Waiting]: '#2e7d32',
  [GameStatus.Active]: '#f57c00',
  [GameStatus.Finished]: '#c62828',
}

const STATUS_LABELS: Record<GameStatus, string> = {
  [GameStatus.Waiting]: 'Waiting for opponent',
  [GameStatus.Active]: 'In progress',
  [GameStatus.Finished]: 'Finished',
}

const COLOR_OPTIONS = [
  { value: 'white' as const, label: '♔' },
  { value: 'random' as const, label: '🎲' },
  { value: 'black' as const, label: '♚' },
]

const RESUME_LABEL = 'Resume'
const VIEW_LABEL = 'View'

export function Lobby() {
  const { user } = getAuth()
  const [waitingGames, setWaitingGames] = useState<GameResponse[]>([])
  const [myGames, setMyGames] = useState<GameResponse[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const navigate = useNavigate()

  const [showCreator, setShowCreator] = useState(false)
  const [selectedCatIdx, setSelectedCatIdx] = useState(1)
  const [selectedPresetIdx, setSelectedPresetIdx] = useState(0)
  const [colorPref, setColorPref] = useState<'white' | 'black' | 'random'>('white')
  const [openCategory, setOpenCategory] = useState(1)

  const modalRef = useRef<HTMLDivElement>(null)

  const loadGames = useCallback(async () => {
    try {
      const [waiting, my] = await Promise.all([getWaitingGames(), getMyGames()])
      setWaitingGames(waiting)
      setMyGames(my)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load games')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadGames()
  }, [loadGames])

  useEffect(() => {
    if (!showCreator) return
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setShowCreator(false)
    }
    window.addEventListener('keydown', handleKey)
    return () => window.removeEventListener('keydown', handleKey)
  }, [showCreator])

  const handleOverlayClick = (e: React.MouseEvent) => {
    if (modalRef.current && !modalRef.current.contains(e.target as Node)) {
      setShowCreator(false)
    }
  }

  const handleCreate = async () => {
    setError('')
    try {
      const cat = ALL_CATEGORIES[selectedCatIdx]
      const preset = cat.presets[selectedPresetIdx]
      const req: CreateGameRequest = {
        time_base: preset.base,
        time_increment: preset.increment,
        color_pref: colorPref,
      }
      const game = await createGame(req)
      setShowCreator(false)
      navigate(`/game/${game.id}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create game')
    }
  }

  const handleJoin = async (id: string) => {
    setError('')
    try {
      await joinGame(id)
      navigate(`/game/${id}`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to join game')
    }
  }

  const handleLogout = () => {
    clearAuth()
    navigate('/login')
  }

  const handleSelectPreset = (catIdx: number, presetIdx: number) => {
    setSelectedCatIdx(catIdx)
    setSelectedPresetIdx(presetIdx)
  }

  const gameButtonLabel = (status: GameStatus): string =>
    status === GameStatus.Active ? RESUME_LABEL : VIEW_LABEL

  const selectedCat = ALL_CATEGORIES[selectedCatIdx]
  const selectedPreset = selectedCat.presets[selectedPresetIdx]
  const selectionLabel = `${selectedCat.icon} ${selectedPreset.label}`

  return (
    <S.LobbyWrapper>
      {showCreator && (
        <S.ModalOverlay onClick={handleOverlayClick}>
          <S.ModalContent ref={modalRef}>
            <S.ModalHeader>
              <h2>New Game</h2>
              <S.ModalClose onClick={() => setShowCreator(false)}>✕</S.ModalClose>
            </S.ModalHeader>

            <S.CreatorSelection>{selectionLabel}</S.CreatorSelection>

            <S.ColorPickerRow>
              {COLOR_OPTIONS.map(opt => (
                <S.ColorPill
                  key={opt.value}
                  $active={colorPref === opt.value}
                  onClick={() => setColorPref(opt.value)}
                >
                  {opt.label}
                </S.ColorPill>
              ))}
            </S.ColorPickerRow>

            <S.TimeSections>
              {ALL_CATEGORIES.map((cat, catIdx) => {
                const isOpen = openCategory === catIdx
                const isCurrentCategory = catIdx === selectedCatIdx
                return (
                  <S.TimeSection key={cat.id}>
                    <S.TimeSectionHeader onClick={() => setOpenCategory(isOpen ? -1 : catIdx)}>
                      <S.TimeSectionLabel>
                        <span>{cat.icon}</span>
                        {cat.label}
                      </S.TimeSectionLabel>
                      <S.Chevron>{isOpen ? '▲' : '▼'}</S.Chevron>
                    </S.TimeSectionHeader>
                    {isOpen && (
                      <S.TimeSectionPills>
                        {cat.presets.map((preset, presetIdx) => {
                          const isSelected = isCurrentCategory && presetIdx === selectedPresetIdx
                          return (
                            <S.TimePill
                              key={preset.label}
                              $selected={isSelected}
                              onClick={() => handleSelectPreset(catIdx, presetIdx)}
                            >
                              {preset.label}
                            </S.TimePill>
                          )
                        })}
                      </S.TimeSectionPills>
                    )}
                  </S.TimeSection>
                )
              })}
            </S.TimeSections>

            <PrimaryButton as={S.CreateButton} onClick={handleCreate}>
              Create Game
            </PrimaryButton>
          </S.ModalContent>
        </S.ModalOverlay>
      )}

      <S.UserBar>
        <h1>Lobby</h1>
        <div>
          <S.UserGreeting>Hello, {user?.username}</S.UserGreeting>
          <SecondaryButton onClick={handleLogout}>Logout</SecondaryButton>
        </div>
      </S.UserBar>

      <S.Actions>
        <PrimaryButton onClick={() => setShowCreator(true)}>+ New Game</PrimaryButton>
        <SecondaryButton onClick={loadGames}>Refresh</SecondaryButton>
      </S.Actions>

      {error && <S.ErrorInline>{error}</S.ErrorInline>}

      {loading ? (
        <p>Loading...</p>
      ) : (
        <>
          <S.SectionHeading>Waiting Games</S.SectionHeading>
          {waitingGames.length === 0 ? (
            <S.EmptyState>
              <p>No games waiting. Create one to get started!</p>
            </S.EmptyState>
          ) : (
            <S.GameList>
              {waitingGames.map(g => (
                <S.GameCard key={g.id}>
                  <S.GameInfo>
                    <span>
                      Host: <strong>{g.white_player.username}</strong>
                    </span>
                    {g.time_control.base > 0 && (
                      <S.TimeBadge>
                        {g.time_control.base / 60}+{g.time_control.increment}
                      </S.TimeBadge>
                    )}
                    <span>
                      Created: <strong>{new Date(g.created_at).toLocaleTimeString()}</strong>
                    </span>
                  </S.GameInfo>
                  <PrimaryButton onClick={() => handleJoin(g.id)}>Join</PrimaryButton>
                </S.GameCard>
              ))}
            </S.GameList>
          )}

          <S.SectionHeading>My Games</S.SectionHeading>
          {myGames.length === 0 ? (
            <S.TextMuted>You have no games yet.</S.TextMuted>
          ) : (
            <S.GameList>
              {myGames.map(g => (
                <S.GameCard key={g.id}>
                  <S.GameInfo>
                    <S.StatusBadge style={{ background: STATUS_COLORS[g.status] }}>
                      {STATUS_LABELS[g.status]}
                    </S.StatusBadge>
                    {g.time_control.base > 0 && (
                      <S.TimeBadge>
                        {g.time_control.base / 60}+{g.time_control.increment}
                      </S.TimeBadge>
                    )}
                  </S.GameInfo>
                  <SecondaryButton onClick={() => navigate(`/game/${g.id}`)}>
                    {gameButtonLabel(g.status)}
                  </SecondaryButton>
                </S.GameCard>
              ))}
            </S.GameList>
          )}
        </>
      )}
    </S.LobbyWrapper>
  )
}
