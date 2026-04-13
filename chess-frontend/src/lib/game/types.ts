export enum GameStatus {
  Waiting = 'waiting',
  Active = 'active',
  Finished = 'finished',
}

export enum GameResult {
  WhiteWins = 'white_wins',
  BlackWins = 'black_wins',
  Draw = 'draw',
}

export enum GameEndReason {
  Checkmate = 'checkmate',
  Stalemate = 'stalemate',
  Resign = 'resign',
  Timeout = 'timeout',
  DrawAgreement = 'draw_agreement',
}

export const RESULT_LABELS: Record<GameResult, string> = {
  [GameResult.WhiteWins]: 'White wins',
  [GameResult.BlackWins]: 'Black wins',
  [GameResult.Draw]: 'Draw',
}

export const END_REASON_LABELS: Record<GameEndReason, string> = {
  [GameEndReason.Checkmate]: 'by checkmate',
  [GameEndReason.Stalemate]: 'by stalemate',
  [GameEndReason.Resign]: 'by resignation',
  [GameEndReason.Timeout]: 'on time',
  [GameEndReason.DrawAgreement]: 'by agreement',
}
