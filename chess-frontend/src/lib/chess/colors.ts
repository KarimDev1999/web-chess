export enum Color {
  White = 'white',
  Black = 'black',
}

export const OPPONENT: Record<Color, Color> = {
  [Color.White]: Color.Black,
  [Color.Black]: Color.White,
}

export const PAWN_START_ROW: Record<Color, number> = {
  [Color.White]: 6,
  [Color.Black]: 1,
}

export const PROMOTION_RANK: Record<Color, number> = {
  [Color.White]: 0,
  [Color.Black]: 7,
}
