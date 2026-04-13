import { getAuth } from '../store/authStore'
import { GameStatus, GameResult, GameEndReason } from '../lib/game'
import { Color } from '../lib/chess'
import { CONTENT_TYPE_JSON } from '../lib/config/http'

export interface PlayerInfo {
  id: string
  username: string
}

export interface GameResponse {
  id: string
  white_player: PlayerInfo
  black_player: PlayerInfo | null
  status: GameStatus
  result: GameResult | null
  end_reason: GameEndReason | null
  fen: string
  turn: Color
  moves: MoveResponse[]
  created_at: string
  updated_at: string
  draw_offer: DrawOfferResponse | null
  time_control: TimeControl
  white_remaining: number | null
  black_remaining: number | null
  last_move_at: string | null
}

export interface TimeControl {
  base: number
  increment: number
}

export interface TimeControlPresets {
  [category: string]: TimeControlPreset[]
}

export interface TimeControlPreset {
  label: string
  base: number
  increment: number
}

export interface CreateGameRequest {
  time_base: number
  time_increment: number
  color_pref: 'white' | 'black' | 'random'
}

export interface MoveResponse {
  from: string
  to: string
  promotion: string | null
  made_at: string
  castle: boolean
  en_passant: boolean
}

export interface MoveRequest {
  from: string
  to: string
}

export interface DrawOfferResponse {
  offered_by: string
  offered_at: string
}

export interface AuthResponse {
  token: string
  user: { id: string; email: string; username: string }
}

const ENDPOINTS = {
  register: '/api/register',
  login: '/api/login',
  timeControls: '/api/time-controls',
  games: '/api/games',
  waitingGames: '/api/games/waiting',
  joinGame: (id: string) => `/api/games/${id}/join`,
  makeMove: (id: string) => `/api/games/${id}/move`,
  getGame: (id: string) => `/api/games/${id}`,
  getMoveHistory: (id: string) => `/api/games/${id}/moves`,
  resign: (id: string) => `/api/games/${id}/resign`,
  offerDraw: (id: string) => `/api/games/${id}/draw/offer`,
  acceptDraw: (id: string) => `/api/games/${id}/draw/accept`,
  declineDraw: (id: string) => `/api/games/${id}/draw/decline`,
} as const

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const { token } = getAuth()
  const headers = new Headers(init?.headers)
  headers.set('Content-Type', CONTENT_TYPE_JSON)
  if (token) headers.set('Authorization', `Bearer ${token}`)

  const res = await fetch(path, { ...init, headers })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `HTTP ${res.status}`)
  }
  if (res.status === 204) return undefined as T

  const ct = res.headers.get('Content-Type')
  if (ct?.includes(CONTENT_TYPE_JSON)) return res.json()
  return undefined as T
}

export async function register(
  email: string,
  password: string,
  username: string,
): Promise<AuthResponse> {
  return request<AuthResponse>(ENDPOINTS.register, {
    method: 'POST',
    body: JSON.stringify({ email, password, username }),
  })
}

export async function login(email: string, password: string): Promise<AuthResponse> {
  return request<AuthResponse>(ENDPOINTS.login, {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export async function createGame(req?: CreateGameRequest): Promise<GameResponse> {
  return request<GameResponse>(ENDPOINTS.games, {
    method: 'POST',
    body: req ? JSON.stringify(req) : undefined,
  })
}

export async function getTimeControlPresets(): Promise<TimeControlPresets> {
  return request<TimeControlPresets>(ENDPOINTS.timeControls)
}

export async function getWaitingGames(): Promise<GameResponse[]> {
  return request<GameResponse[]>(ENDPOINTS.waitingGames)
}

export async function getMyGames(): Promise<GameResponse[]> {
  return request<GameResponse[]>(ENDPOINTS.games)
}

export async function getGame(id: string): Promise<GameResponse> {
  return request<GameResponse>(ENDPOINTS.getGame(id))
}

export async function joinGame(id: string): Promise<GameResponse> {
  return request<GameResponse>(ENDPOINTS.joinGame(id), { method: 'POST' })
}

export async function makeMove(id: string, move: MoveRequest): Promise<GameResponse> {
  return request<GameResponse>(ENDPOINTS.makeMove(id), {
    method: 'POST',
    body: JSON.stringify(move),
  })
}

export async function getMoveHistory(id: string): Promise<MoveResponse[]> {
  return request<MoveResponse[]>(ENDPOINTS.getMoveHistory(id))
}

export async function resignGame(id: string): Promise<GameResponse> {
  return request<GameResponse>(ENDPOINTS.resign(id), { method: 'POST' })
}

export async function offerDraw(id: string): Promise<GameResponse> {
  return request<GameResponse>(ENDPOINTS.offerDraw(id), { method: 'POST' })
}

export async function acceptDraw(id: string): Promise<GameResponse> {
  return request<GameResponse>(ENDPOINTS.acceptDraw(id), { method: 'POST' })
}

export async function declineDraw(id: string): Promise<GameResponse> {
  return request<GameResponse>(ENDPOINTS.declineDraw(id), { method: 'POST' })
}
