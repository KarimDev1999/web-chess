import { useState, useEffect, useCallback, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { getAuth, clearAuth } from '../store/authStore'
import {
  getWaitingGames,
  getMyGames,
  createGame,
  joinGame,
  type GameResponse,
  type CreateGameRequest,
} from '../api/chessApi'
import { GameStatus } from '../lib/game'
import { ALL_CATEGORIES } from '../lib/game/timeControls'

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
    <div className="lobby">
      {/* Modal overlay */}
      {showCreator && (
        <div className="modal-overlay" onClick={handleOverlayClick}>
          <div className="modal-content" ref={modalRef}>
            <div className="modal-header">
              <h2>New Game</h2>
              <button className="modal-close" onClick={() => setShowCreator(false)}>
                ✕
              </button>
            </div>

            {/* Current selection */}
            <div className="creator-selection">{selectionLabel}</div>

            {/* Color picker */}
            <div className="color-picker-row">
              {COLOR_OPTIONS.map(opt => (
                <button
                  key={opt.value}
                  className={`color-pill ${colorPref === opt.value ? 'active' : ''}`}
                  onClick={() => setColorPref(opt.value)}
                >
                  {opt.label}
                </button>
              ))}
            </div>

            {/* Time control sections */}
            <div className="time-sections">
              {ALL_CATEGORIES.map((cat, catIdx) => {
                const isOpen = openCategory === catIdx
                const isCurrentCategory = catIdx === selectedCatIdx
                return (
                  <div key={cat.id} className="time-section">
                    <button
                      className="time-section-header"
                      onClick={() => setOpenCategory(isOpen ? -1 : catIdx)}
                    >
                      <span className="time-section-label">
                        <span className="section-icon">{cat.icon}</span>
                        {cat.label}
                      </span>
                      <span className="chevron">{isOpen ? '▲' : '▼'}</span>
                    </button>
                    {isOpen && (
                      <div className="time-section-pills">
                        {cat.presets.map((preset, presetIdx) => {
                          const isSelected = isCurrentCategory && presetIdx === selectedPresetIdx
                          return (
                            <button
                              key={preset.label}
                              className={`time-pill ${isSelected ? 'selected' : ''}`}
                              onClick={() => handleSelectPreset(catIdx, presetIdx)}
                            >
                              {preset.label}
                            </button>
                          )
                        })}
                      </div>
                    )}
                  </div>
                )
              })}
            </div>

            <button className="btn-primary btn-create-game" onClick={handleCreate}>
              Create Game
            </button>
          </div>
        </div>
      )}

      <div className="user-bar">
        <h1>Lobby</h1>
        <div>
          <span className="user-greeting">Hello, {user?.username}</span>
          <button className="btn-secondary" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </div>

      <div className="actions">
        <button className="btn-primary" onClick={() => setShowCreator(true)}>
          + New Game
        </button>
        <button className="btn-secondary" onClick={loadGames}>
          Refresh
        </button>
      </div>

      {error && <div className="error-inline">{error}</div>}

      {loading ? (
        <p>Loading...</p>
      ) : (
        <>
          <h2 className="section-heading">Waiting Games</h2>
          {waitingGames.length === 0 ? (
            <div className="empty-state section-gap">
              <p>No games waiting. Create one to get started!</p>
            </div>
          ) : (
            <div className="game-list section-gap">
              {waitingGames.map(g => (
                <div key={g.id} className="game-card">
                  <div className="info">
                    <span>
                      Host: <strong>{g.white_player.username}</strong>
                    </span>
                    {g.time_control.base > 0 && (
                      <span className="time-badge">
                        {g.time_control.base / 60}+{g.time_control.increment}
                      </span>
                    )}
                    <span>
                      Created: <strong>{new Date(g.created_at).toLocaleTimeString()}</strong>
                    </span>
                  </div>
                  <button className="btn-primary" onClick={() => handleJoin(g.id)}>
                    Join
                  </button>
                </div>
              ))}
            </div>
          )}

          <h2 className="section-heading">My Games</h2>
          {myGames.length === 0 ? (
            <p className="text-muted">You have no games yet.</p>
          ) : (
            <div className="game-list">
              {myGames.map(g => (
                <div key={g.id} className="game-card">
                  <div className="info">
                    <span className="status-badge" style={{ background: STATUS_COLORS[g.status] }}>
                      {STATUS_LABELS[g.status]}
                    </span>
                    {g.time_control.base > 0 && (
                      <span className="time-badge">
                        {g.time_control.base / 60}+{g.time_control.increment}
                      </span>
                    )}
                  </div>
                  <button className="btn-secondary" onClick={() => navigate(`/game/${g.id}`)}>
                    {gameButtonLabel(g.status)}
                  </button>
                </div>
              ))}
            </div>
          )}
        </>
      )}
    </div>
  )
}
