import styled, { keyframes, css } from 'styled-components'
import { colors } from '../../styles/styled'

export const GameViewWrapper = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px 20px;
`

export const GameContent = styled.div`
  display: flex;
  gap: 24px;
  align-items: flex-start;
`

export const BoardColumn = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
`

export const BoardContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  width: 100%;
  max-width: 560px;
`

export const GameHeaderBar = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  width: 100%;
  max-width: 560px;
`

export const GameStatus = styled.span<{ $status: string }>`
  padding: 8px 16px;
  border-radius: 6px;
  font-weight: 600;
  font-size: 14px;
  background: ${props => props.$status === 'active' ? colors.success : colors.danger};
`

export const WsIndicator = styled.span`
  font-size: 13px;
  color: ${colors.textMuted};
`

export const ErrorBanner = styled.div`
  color: ${colors.red};
  padding: 8px 12px;
  background: ${colors.redBg};
  border-radius: 6px;
  margin-bottom: 12px;
  font-size: 14px;
  width: 100%;
  max-width: 560px;
`

export const GameSidebar = styled.div`
  width: 280px;
  flex-shrink: 0;
`

export const MoveHistory = styled.div`
  background: ${colors.bgCard};
  border-radius: 8px;
  padding: 16px;
  max-height: 600px;
  overflow-y: auto;

  h3 {
    margin-bottom: 12px;
    font-size: 16px;
  }
`

export const MovesList = styled.div`
  display: grid;
  grid-template-columns: 30px 1fr 1fr;
  gap: 4px 8px;
  font-size: 14px;

  .move {
    padding: 2px 6px;
    border-radius: 3px;
    cursor: default;

    &:hover {
      background: ${colors.bgHover};
    }
  }

  .move-num {
    color: ${colors.textMuted};
  }
`

export const MoveRow = styled.div`
  display: contents;
`

export const NoMoves = styled.p`
  color: ${colors.textMuted};
  font-size: 14px;
`

export const WaitingNotice = styled.div`
  margin-top: 16px;
  background: ${colors.bgSection};
  border-radius: 8px;
  padding: 16px;
  text-align: center;

  p {
    margin: 0;
    font-size: 14px;
    color: ${colors.textSecondary};
  }
`

export const GameOverNotice = styled.div`
  margin-top: 16px;
  background: ${colors.bgSection};
  border-radius: 8px;
  padding: 16px;
  text-align: center;

  p {
    margin: 0;
    font-size: 14px;
    color: ${colors.textSecondary};
  }
`

export const DrawOfferInline = styled.div`
  margin-top: 16px;
  background: ${colors.bgSection};
  border-radius: 8px;
  padding: 16px;
  text-align: center;

  .draw-offer-text {
    margin: 0 0 12px;
    font-size: 14px;
    color: ${colors.text};
  }

  .draw-offer-timer {
    color: ${colors.red};
    font-weight: 600;
  }
`

export const DrawOfferButtons = styled.div`
  display: flex;
  gap: 12px;
  justify-content: center;
`

export const GameActions = styled.div`
  margin-top: 16px;
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
`

export const DrawPendingLabel = styled.span`
  font-size: 13px;
  color: ${colors.textMuted};
  font-style: italic;
`

const clockPulse = keyframes`
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
`

export const PlayerClock = styled.span<{ $ticking?: boolean }>`
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
    animation: ${clockPulse} 1s ease-in-out infinite;
  `}
`
