import { useState, FormEvent } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { login } from '../../api/chessApi'
import { setAuth } from '../../store/authStore'
import * as S from './styles'
import { PrimaryButton, TextInput } from '../../styles/styled'

export function Login() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const data = await login(email, password)
      setAuth(data.user, data.token)
      navigate('/lobby')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <S.Container>
      <S.Card>
        <h1>Web Chess</h1>
        <S.Subtitle>Sign in to your account</S.Subtitle>
        <S.Form onSubmit={handleSubmit}>
          {error && <S.ErrorText>{error}</S.ErrorText>}
          <S.FormGroup>
            <label>Email</label>
            <TextInput type="email" value={email} onChange={e => setEmail(e.target.value)} required />
          </S.FormGroup>
          <S.FormGroup>
            <label>Password</label>
            <TextInput type="password" value={password} onChange={e => setPassword(e.target.value)} required />
          </S.FormGroup>
          <PrimaryButton as={S.SubmitButton} type="submit" disabled={loading}>
            {loading ? 'Logging in...' : 'Sign In'}
          </PrimaryButton>
        </S.Form>
        <S.LinkText>
          Don&apos;t have an account? <Link to="/register">Create one</Link>
        </S.LinkText>
      </S.Card>
    </S.Container>
  )
}
