import styled from 'styled-components'

const colors = {
  bg: '#1a1a2e',
  bgCard: '#16213e',
  bgInput: '#16213e',
  bgHover: '#1e3a5f',
  bgSection: '#1a2744',
  bgTimeBadge: '#1e3a5f',
  text: '#e0e0e0',
  textMuted: '#78909c',
  textSubtle: '#607d8b',
  textSecondary: '#90a4ae',
  accent: '#4fc3f7',
  accentHover: '#29b6f6',
  secondaryBg: '#37474f',
  secondaryHover: '#455a64',
  danger: '#c62828',
  dangerHover: '#b71c1c',
  success: '#2e7d32',
  successDark: '#1b5e20',
  successLight: '#a5d6a7',
  green: '#4caf50',
  orange: '#f57c00',
  red: '#ef5350',
  redBg: '#2d1b1b',
  boardLight: '#f0d9b5',
  boardDark: '#b58863',
  selectedSquare: '#7fc97f',
  lastMoveSquare: '#cdd26a',
  dragSourceSquare: '#a9a9a9',
  error: '#ef5350',
}

export const Button = styled.button`
  cursor: pointer;
  border: none;
  border-radius: 6px;
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.2s;

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`

export const PrimaryButton = styled(Button)`
  background: ${colors.accent};
  color: ${colors.bg};

  &:hover:not(:disabled) {
    background: ${colors.accentHover};
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(79, 195, 247, 0.3);
  }

  &:active:not(:disabled) {
    transform: translateY(0);
  }
`

export const SecondaryButton = styled(Button)`
  background: ${colors.secondaryBg};
  color: ${colors.text};

  &:hover:not(:disabled) {
    background: ${colors.secondaryHover};
  }
`

export const DangerButton = styled(Button)`
  background: ${colors.danger};
  color: ${colors.text};

  &:hover:not(:disabled) {
    background: ${colors.dangerHover};
  }
`

export const TextInput = styled.input`
  padding: 12px 16px;
  border-radius: 8px;
  border: 1px solid ${colors.secondaryBg};
  background: ${colors.bgInput};
  color: ${colors.text};
  font-size: 16px;
  width: 100%;
  transition: all 0.2s;

  &:focus {
    outline: none;
    border-color: ${colors.accent};
    box-shadow: 0 0 0 3px rgba(79, 195, 247, 0.15);
  }
`

export { colors }
