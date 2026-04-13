import { Color, OPPONENT, PAWN_START_ROW } from './colors'
import { PieceType, FEN_CHAR_TO_TYPE } from './pieces'
import {
  BOARD_SIZE,
  FILE_LETTERS,
  KNIGHT_OFFSETS,
  DIAGONAL_DIRS,
  ORTHOGONAL_DIRS,
  KING_DIRS,
} from './moves'

export interface Piece {
  type: PieceType
  color: Color
}

export type Board = (Piece | null)[][]

export interface Pos {
  row: number
  col: number
}

/** Castling availability encoded as a string like "KQkq", "kq", "-", etc. */
export type CastlingRights = string

export interface FENMeta {
  castling: CastlingRights
  enPassantTarget: Pos | null
}

export function inBounds(row: number, col: number): boolean {
  return row >= 0 && row < BOARD_SIZE && col >= 0 && col < BOARD_SIZE
}

export function toAlgebraic(row: number, col: number): string {
  return `${FILE_LETTERS[col]}${BOARD_SIZE - row}`
}

export function fromAlgebraic(s: string): Pos {
  return { row: BOARD_SIZE - parseInt(s[1], 10), col: s.charCodeAt(0) - 97 }
}

/** Parse the board placement part of a FEN string (just pieces). */
export function parseFEN(fen: string): Board {
  const [placement] = fen.split(' ')
  const board: Board = []

  for (const row of placement.split('/')) {
    const cells: (Piece | null)[] = []
    for (const ch of row) {
      const n = parseInt(ch, 10)
      if (!isNaN(n)) {
        for (let i = 0; i < n; i++) cells.push(null)
      } else {
        const color: Color = ch === ch.toUpperCase() ? Color.White : Color.Black
        cells.push({ type: FEN_CHAR_TO_TYPE[ch.toLowerCase()], color })
      }
    }
    board.push(cells)
  }
  return board
}

/** Parse castling rights and en-passant target from a full FEN string. */
export function parseFENMeta(fen: string): FENMeta {
  const parts = fen.split(' ')
  const castling = parts.length >= 3 ? parts[2] : '-'
  const enPassantTarget: Pos | null =
    parts.length >= 4 && parts[3] !== '-' ? fromAlgebraic(parts[3]) : null
  return { castling, enPassantTarget }
}

export function getPseudoLegalMoves(
  board: Board,
  row: number,
  col: number,
  turn: Color,
  fenMeta?: FENMeta,
): Pos[] {
  const piece = board[row]?.[col]
  if (!piece || piece.color !== turn) return []

  const moves: Pos[] = []
  const enemy = OPPONENT[turn]

  const addIf = (r: number, c: number) => {
    if (!inBounds(r, c)) return false
    const target = board[r][c]
    if (!target) {
      moves.push({ row: r, col: c })
      return true
    }
    if (target.color === enemy) {
      moves.push({ row: r, col: c })
      return false
    }
    return false
  }

  const slide = (dr: number, dc: number) => {
    for (let i = 1; i < BOARD_SIZE; i++) {
      if (!addIf(row + dr * i, col + dc * i)) break
    }
  }

  switch (piece.type) {
    case PieceType.Pawn: {
      const dir = piece.color === Color.White ? -1 : 1
      const startRow = PAWN_START_ROW[piece.color]

      if (!board[row + dir]?.[col]) {
        moves.push({ row: row + dir, col })
        if (row === startRow && !board[row + 2 * dir]?.[col])
          moves.push({ row: row + 2 * dir, col })
      }
      for (const dc of [-1, 1]) {
        const tr = row + dir,
          tc = col + dc
        if (inBounds(tr, tc)) {
          const t = board[tr][tc]
          if (t && t.color === enemy) moves.push({ row: tr, col: tc })
          // En-passant capture
          if (
            fenMeta?.enPassantTarget &&
            fenMeta.enPassantTarget.row === tr &&
            fenMeta.enPassantTarget.col === tc
          ) {
            moves.push({ row: tr, col: tc })
          }
        }
      }
      break
    }
    case PieceType.Knight:
      for (const [dr, dc] of KNIGHT_OFFSETS) addIf(row + dr, col + dc)
      break
    case PieceType.Bishop:
      for (const [dr, dc] of DIAGONAL_DIRS) slide(dr, dc)
      break
    case PieceType.Rook:
      for (const [dr, dc] of ORTHOGONAL_DIRS) slide(dr, dc)
      break
    case PieceType.Queen:
      for (const [dr, dc] of [...DIAGONAL_DIRS, ...ORTHOGONAL_DIRS]) slide(dr, dc)
      break
    case PieceType.King: {
      for (const [dr, dc] of KING_DIRS) addIf(row + dr, col + dc)
      // Castling moves
      if (fenMeta && fenMeta.castling !== '-') {
        const isWhite = piece.color === Color.White
        const homeRow = isWhite ? 7 : 0
        if (row === homeRow) {
          const canK = isWhite ? fenMeta.castling.includes('K') : fenMeta.castling.includes('k')
          const canQ = isWhite ? fenMeta.castling.includes('Q') : fenMeta.castling.includes('q')
          if (canK && !board[homeRow][5] && !board[homeRow][6]) {
            moves.push({ row: homeRow, col: 6 })
          }
          if (canQ && !board[homeRow][3] && !board[homeRow][2] && !board[homeRow][1]) {
            moves.push({ row: homeRow, col: 2 })
          }
        }
      }
      break
    }
  }
  return moves
}

/** Find the king position for the given color. */
export function findKing(board: Board, color: Color): Pos | null {
  for (let r = 0; r < BOARD_SIZE; r++) {
    for (let c = 0; c < BOARD_SIZE; c++) {
      const p = board[r]?.[c]
      if (p?.type === PieceType.King && p.color === color) {
        return { row: r, col: c }
      }
    }
  }
  return null
}

/** Check if the given position is attacked by any piece of the opponent color. */
function isSquareAttackedBy(board: Board, row: number, col: number, byColor: Color): boolean {
  const enemy = byColor

  // Knight attacks
  for (const [dr, dc] of KNIGHT_OFFSETS) {
    const r = row + dr,
      c = col + dc
    if (inBounds(r, c)) {
      const p = board[r][c]
      if (p?.color === enemy && p.type === PieceType.Knight) return true
    }
  }

  // Pawn attacks
  const pawnDir = enemy === Color.White ? 1 : -1
  for (const dc of [-1, 1]) {
    const r = row + pawnDir,
      c = col + dc
    if (inBounds(r, c)) {
      const p = board[r][c]
      if (p?.color === enemy && p.type === PieceType.Pawn) return true
    }
  }

  // King attacks
  for (const [dr, dc] of KING_DIRS) {
    const r = row + dr,
      c = col + dc
    if (inBounds(r, c)) {
      const p = board[r][c]
      if (p?.color === enemy && p.type === PieceType.King) return true
    }
  }

  // Sliding attacks (queen, rook, bishop)
  const trySlide = (dr: number, dc: number, types: PieceType[]): boolean => {
    for (let i = 1; i < BOARD_SIZE; i++) {
      const r = row + dr * i,
        c = col + dc * i
      if (!inBounds(r, c)) break
      const p = board[r][c]
      if (p) {
        if (p.color === enemy && types.includes(p.type)) return true
        break
      }
    }
    return false
  }

  // Orthogonal (rook/queen)
  for (const [dr, dc] of ORTHOGONAL_DIRS) {
    if (trySlide(dr, dc, [PieceType.Rook, PieceType.Queen])) return true
  }

  // Diagonal (bishop/queen)
  for (const [dr, dc] of DIAGONAL_DIRS) {
    if (trySlide(dr, dc, [PieceType.Bishop, PieceType.Queen])) return true
  }

  return false
}

/** Return true if the player whose turn it is is currently in check. */
export function isInCheck(board: Board, turn: Color): boolean {
  const king = findKing(board, turn)
  if (!king) return false
  return isSquareAttackedBy(board, king.row, king.col, OPPONENT[turn])
}
