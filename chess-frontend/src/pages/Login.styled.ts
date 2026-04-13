import styled from 'styled-components'
import { colors } from '../styles/styled'

export const Container = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #0d1b2a 0%, #1a1a2e 50%, #16213e 100%);
`

export const Card = styled.div`
  background: ${colors.bgCard};
  padding: 48px 40px;
  border-radius: 16px;
  width: 100%;
  max-width: 420px;
  box-shadow:
    0 4px 24px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(79, 195, 247, 0.08);

  h1 {
    margin-bottom: 8px;
    text-align: center;
    font-size: 28px;
    font-weight: 700;
    color: ${colors.accent};
    letter-spacing: -0.5px;
  }
`

export const Subtitle = styled.p`
  text-align: center;
  color: ${colors.textSecondary};
  font-size: 14px;
  margin-bottom: 32px;
`

export const Form = styled.form`
  display: flex;
  flex-direction: column;
`

export const FormGroup = styled.div`
  margin-bottom: 20px;

  label {
    display: block;
    margin-bottom: 8px;
    font-size: 13px;
    font-weight: 600;
    color: ${colors.textSecondary};
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
`

export const SubmitButton = styled.button`
  width: 100%;
  margin-top: 12px;
  padding: 14px;
  font-size: 16px;
  font-weight: 700;
  border-radius: 8px;
  letter-spacing: 0.3px;
`

export const ErrorText = styled.div`
  color: ${colors.red};
  font-size: 13px;
  margin-bottom: 16px;
  padding: 10px 14px;
  background: rgba(239, 83, 80, 0.1);
  border-radius: 8px;
  border-left: 3px solid ${colors.red};
`

export const LinkText = styled.p`
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: ${colors.textSecondary};

  a {
    color: ${colors.accent};
    font-weight: 600;
    text-decoration: none;
    transition: color 0.2s;

    &:hover {
      color: ${colors.accentHover};
      text-decoration: underline;
    }
  }
`
