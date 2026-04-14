import styled, { keyframes } from 'styled-components'
import { colors } from '../../styles/styled'

export const LobbyWrapper = styled.div`
  max-width: 900px;
  margin: 0 auto;
  padding: 40px 20px;
`

export const UserBar = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  h1 {
    margin-bottom: 0;
  }
`

export const UserGreeting = styled.span`
  margin-right: 16px;
`

export const Actions = styled.div`
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
`

export const SectionHeading = styled.h2`
  margin-bottom: 12px;
  font-size: 18px;
`

export const SectionGap = styled.div`
  margin-bottom: 32px;
`

export const ErrorInline = styled.div`
  color: ${colors.red};
  margin-bottom: 16px;
`

export const TextMuted = styled.p`
  color: ${colors.textMuted};
`

export const GameList = styled.div`
  display: grid;
  gap: 12px;
`

export const GameCard = styled.div`
  background: ${colors.bgCard};
  padding: 16px 20px;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
`

export const GameInfo = styled.div`
  display: flex;
  gap: 16px;
  align-items: center;
  font-size: 14px;
  color: ${colors.textSecondary};

  strong {
    color: ${colors.text};
  }
`

export const TimeBadge = styled.span`
  background: ${colors.bgTimeBadge};
  color: ${colors.accent};
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
`

export const StatusBadge = styled.span`
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
`

export const EmptyState = styled.div`
  text-align: center;
  padding: 40px;
  color: ${colors.textMuted};
`

const fadeIn = keyframes`
  from { opacity: 0; }
  to { opacity: 1; }
`

const slideUp = keyframes`
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
`

export const ModalOverlay = styled.div`
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: ${fadeIn} 0.15s ease-out;
`

export const ModalContent = styled.div`
  background: ${colors.bgCard};
  border-radius: 12px;
  padding: 24px;
  width: 90%;
  max-width: 420px;
  max-height: 85vh;
  overflow-y: auto;
  animation: ${slideUp} 0.2s ease-out;
`

export const ModalHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  h2 {
    margin: 0;
    font-size: 20px;
    color: ${colors.text};
  }
`

export const ModalClose = styled.button`
  background: none;
  border: none;
  color: ${colors.textMuted};
  font-size: 20px;
  cursor: pointer;
  padding: 4px 8px;
  line-height: 1;

  &:hover {
    color: ${colors.text};
  }
`

export const CreatorSelection = styled.div`
  text-align: center;
  font-size: 18px;
  font-weight: 700;
  color: ${colors.accent};
  background: ${colors.bgHover};
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 16px;
`

export const ColorPickerRow = styled.div`
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 20px;
`

export const ColorPill = styled.button<{ $active?: boolean }>`
  background: ${colors.secondaryBg};
  color: ${colors.text};
  border: 2px solid transparent;
  border-radius: 8px;
  padding: 8px 16px;
  font-size: 18px;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;

  &:hover {
    background: ${colors.secondaryHover};
  }

  ${props => props.$active && `
    border-color: ${colors.green};
    background: ${colors.successDark};
  `}
`

export const TimeSections = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 20px;
`

export const TimeSection = styled.div`
  background: ${colors.bgSection};
  border-radius: 8px;
  overflow: hidden;
`

export const TimeSectionHeader = styled.button`
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: none;
  border: none;
  color: ${colors.text};
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s;

  &:hover {
    background: ${colors.bgHover};
  }
`

export const TimeSectionLabel = styled.span`
  display: flex;
  align-items: center;
  gap: 8px;
`

export const Chevron = styled.span`
  font-size: 12px;
  color: ${colors.textMuted};
`

export const TimeSectionPills = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 4px 16px 12px;
  border-top: 1px solid #2a3a5a;
`

export const TimePill = styled.button<{ $selected?: boolean }>`
  background: ${colors.secondaryBg};
  color: ${colors.text};
  border: 2px solid transparent;
  border-radius: 6px;
  padding: 8px 14px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s;

  &:hover {
    background: ${colors.secondaryHover};
  }

  ${props => props.$selected && `
    border-color: ${colors.green};
    background: ${colors.successDark};
    color: ${colors.successLight};
  `}
`

export const CreateButton = styled.button`
  width: 100%;
  padding: 14px;
  font-size: 16px;
  font-weight: 700;
`
