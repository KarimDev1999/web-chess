export enum PieceType {
  King = 'king',
  Queen = 'queen',
  Rook = 'rook',
  Bishop = 'bishop',
  Knight = 'knight',
  Pawn = 'pawn',
}

export const FEN_CHAR_TO_TYPE: Record<string, PieceType> = {
  k: PieceType.King,
  q: PieceType.Queen,
  r: PieceType.Rook,
  b: PieceType.Bishop,
  n: PieceType.Knight,
  p: PieceType.Pawn,
}
