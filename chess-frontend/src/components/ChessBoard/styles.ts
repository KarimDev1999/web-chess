import styled, { keyframes } from 'styled-components'
import { colors } from '../../styles/styled'

export const Board = styled.div<{ $dragging?: boolean }>`
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  grid-template-rows: repeat(8, 1fr);
  width: 100%;
  max-width: 560px;
  aspect-ratio: 1 / 1;
  border: 3px solid ${colors.secondaryBg};
  border-radius: 4px;
  user-select: none;
  touch-action: none;

  ${props => props.$dragging && `
    cursor: grabbing !important;
    * { cursor: grabbing !important; }
  `}
`

interface SquareProps {
  $isLight: boolean
  $selected: boolean
  $dragSource: boolean
  $dragOver: boolean
  $kingCheck: boolean
}

export const Square = styled.div<SquareProps>`
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40px;
  position: relative;
  background: ${props => {
    if (props.$selected) return colors.selectedSquare
    if (props.$dragSource) return colors.dragSourceSquare
    if (props.$dragOver) return colors.lastMoveSquare
    return props.$isLight ? colors.boardLight : colors.boardDark
  }};
  cursor: pointer;

  img {
    filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
    transition: opacity 0.05s;
    pointer-events: none;
  }

  ${props => props.$dragSource && `
    img {
      opacity: 0.3;
    }
  `}
`

const checkPulse = keyframes`
  0%, 100% { box-shadow: inset 0 0 12px 4px rgba(255, 0, 0, 0.5); }
  50% { box-shadow: inset 0 0 20px 8px rgba(255, 0, 0, 0.7); }
`

export const KingCheckSquare = styled(Square)`
  animation: ${checkPulse} 1s ease-in-out infinite;
`

export const MoveIndicator = styled.span`
  position: absolute;
  width: 28%;
  height: 28%;
  background: rgba(0, 0, 0, 0.25);
  border-radius: 50%;
  pointer-events: none;
`

export const CaptureIndicator = styled.span`
  position: absolute;
  width: 90%;
  height: 90%;
  border: 5px solid rgba(0, 0, 0, 0.25);
  border-radius: 50%;
  box-sizing: border-box;
  pointer-events: none;
`

export const DraggedPiece = styled.div`
  position: fixed;
  width: 60px;
  height: 60px;
  pointer-events: none;
  z-index: 9999;
  filter: drop-shadow(0 6px 16px rgba(0, 0, 0, 0.6));
  transform: scale(1.15);

  img {
    width: 100%;
    height: 100%;
  }
`
