interface PlayerProfileProps {
  username: string
  color: 'white' | 'black'
  isActive: boolean
  isPresent: boolean
  capturedPieces?: string[]
  clock?: string
}

const pieceSymbols: Record<string, string> = {
  P: '♙',
  R: '♖',
  N: '♘',
  B: '♗',
  Q: '♕',
  K: '♔',
  p: '♟',
  r: '♜',
  n: '♞',
  b: '♝',
  q: '♛',
  k: '♚',
}

export function PlayerProfile({
  username,
  color,
  isActive,
  isPresent,
  capturedPieces = [],
  clock,
}: PlayerProfileProps) {
  const initial = username.charAt(0).toUpperCase()
  const bgColor = color === 'white' ? '#f0d9b5' : '#b58863'
  const textColor = color === 'white' ? '#1a1a2e' : '#e0e0e0'

  return (
    <div className={`player-profile player-${color} ${isActive ? 'active' : ''}`}>
      <div className="player-info-row">
        <div className="player-avatar" style={{ background: bgColor, color: textColor }}>
          {initial}
        </div>
        <div className="player-details">
          <span className="player-name">{username}</span>
          <span className="player-color-label">{color === 'white' ? 'White' : 'Black'}</span>
        </div>
        {clock && <span className={`player-clock ${isActive ? 'ticking' : ''}`}>{clock}</span>}
        <span
          className={`presence-dot ${isPresent ? 'online' : 'offline'}`}
          title={isPresent ? 'Online' : 'Offline'}
        />
      </div>
      {capturedPieces.length > 0 && (
        <div className="captured-pieces">
          {capturedPieces.map((p, i) => (
            <span key={i} className="captured-piece">
              {pieceSymbols[p] || p}
            </span>
          ))}
        </div>
      )}
    </div>
  )
}
