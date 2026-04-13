import styled, { keyframes, css } from 'styled-components'
import { colors } from '../styles/styled'

const pulse = keyframes`
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
`

interface ProfileProps {
  $isActive: boolean
  $color: 'white' | 'black'
}

export const ProfileWrapper = styled.div<ProfileProps>`
  width: 100%;
  background: ${colors.bgCard};
  border-radius: 8px;
  padding: 10px 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  transition: box-shadow 0.2s;

  ${props => props.$isActive && `
    box-shadow: 0 0 0 2px ${colors.green};
  `}
`

export const InfoRow = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`

export const Avatar = styled.div<{ $bg: string; $text: string }>`
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
  flex-shrink: 0;
  background: ${props => props.$bg};
  color: ${props => props.$text};
`

export const Details = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0;
`

export const Name = styled.span`
  font-weight: 600;
  font-size: 14px;
  color: ${colors.text};
  line-height: 1.2;
`

export const ColorLabel = styled.span`
  font-size: 11px;
  color: ${colors.textMuted};
  line-height: 1.2;
`

export const Clock = styled.span<{ $ticking?: boolean }>`
  margin-left: auto;
  font-family: 'Courier New', monospace;
  font-size: 18px;
  font-weight: 700;
  color: ${colors.text};
  padding: 2px 8px;
  background: ${colors.bgSection};
  border-radius: 6px;
  line-height: 1.2;

  ${props => props.$ticking && css`
    color: ${colors.accent};
    animation: ${pulse} 1s ease-in-out infinite;
  `}
`

export const PresenceDot = styled.span<{ $online: boolean }>`
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-left: auto;
  flex-shrink: 0;
  transition: background 0.2s;
  background: ${props => props.$online ? colors.green : colors.textSubtle};

  ${props => props.$online && css`
    animation: ${pulse} 1.5s ease-in-out infinite;
  `}
`

export const CapturedPieces = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
  padding-left: 46px;
  min-height: 18px;
`

export const CapturedPiece = styled.span`
  font-size: 16px;
  color: ${colors.textSecondary};
`
