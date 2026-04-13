import { useState, FormEvent } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { register } from '../api/chessApi'
import { setAuth } from '../store/authStore'
import * as S from './Login.styled'
import { PrimaryButton, TextInput } from '../styles/styled'

export function Register() {
  const [email, setEmail] = useState('')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const data = await register(email, password, username)
      setAuth(data.user, data.token)
      navigate('/lobby')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <S.Container>
      <S.Card>
        <h1>Create Account</h1>
        <S.Subtitle>Join the chess community</S.Subtitle>
        <S.Form onSubmit={handleSubmit}>
          {error && <S.ErrorText>{error}</S.ErrorText>}
          <S.FormGroup>
            <label>Email</label>
            <TextInput type="email" value={email} onChange={e => setEmail(e.target.value)} required />
          </S.FormGroup>
          <S.FormGroup>
            <label>Username</label>
            <TextInput value={username} onChange={e => setUsername(e.target.value)} required />
          </S.FormGroup>
          <S.FormGroup>
            <label>Password</label>
            <TextInput type="password" value={password} onChange={e => setPassword(e.target.value)} required />
          </S.FormGroup>
          <PrimaryButton as={S.SubmitButton} type="submit" disabled={loading}>
            {loading ? 'Creating account...' : 'Sign Up'}
          </PrimaryButton>
        </S.Form>
        <S.LinkText>
          Already have an account? <Link to="/login">Sign in</Link>
        </S.LinkText>
      </S.Card>
    </S.Container>
  )
}
