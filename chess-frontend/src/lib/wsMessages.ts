export const SERVER_EVENT = {
  GAME_CREATED: 'game.created',
  GAME_JOINED: 'game.joined',
  MOVE_MADE: 'game.move_made',
  DRAW_OFFERED: 'game.draw_offered',
  DRAW_DECLINED: 'game.draw_declined',
  DRAW_OFFER_EXPIRED: 'game.draw_offer_expired',
  GAME_RESIGNED: 'game.resigned',
  DRAW_ACCEPTED: 'game.draw_accepted',
  GAME_TIMED_OUT: 'game.timed_out',
  PRESENCE: 'presence',
} as const

export const CLIENT_TYPE = {
  PING: 'ping',
  JOIN_GAME: 'join_game',
  LEAVE_GAME: 'leave_game',
  RESIGN: 'resign',
  OFFER_DRAW: 'offer_draw',
  ACCEPT_DRAW: 'accept_draw',
  DECLINE_DRAW: 'decline_draw',
} as const

export const KEY = {
  TYPE: 'type',
  GAME_ID: 'game_id',
  PLAYERS: 'players',
  OFFERED_BY: 'offered_by',
  DATA: 'data',
} as const

export interface ServerMessage {
  type: string
  game_id?: string
  [key: string]: unknown
}

export interface ClientMessage {
  type: (typeof CLIENT_TYPE)[keyof typeof CLIENT_TYPE]
  data?: Record<string, unknown>
}
