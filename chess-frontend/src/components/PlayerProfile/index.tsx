import * as S from './styles'
import { colors } from '../../styles/styled'

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
  const bgColor = color === 'white' ? colors.boardLight : colors.boardDark
  const textColor = color === 'white' ? '#1a1a2e' : colors.text

  return (
    <S.ProfileWrapper $isActive={isActive} $color={color}>
      <S.InfoRow>
        <S.Avatar $bg={bgColor} $text={textColor}>{initial}</S.Avatar>
        <S.Details>
          <S.Name>{username}</S.Name>
          <S.ColorLabel>{color === 'white' ? 'White' : 'Black'}</S.ColorLabel>
        </S.Details>
        {clock && <S.Clock $ticking={isActive}>{clock}</S.Clock>}
        <S.PresenceDot $online={isPresent} title={isPresent ? 'Online' : 'Offline'} />
      </S.InfoRow>
      {capturedPieces.length > 0 && (
        <S.CapturedPieces>
          {capturedPieces.map((p, i) => (
            <S.CapturedPiece key={i}>{pieceSymbols[p] || p}</S.CapturedPiece>
          ))}
        </S.CapturedPieces>
      )}
    </S.ProfileWrapper>
  )
}
