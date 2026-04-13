import { Color } from '../chess/colors'

export const PIECE_SVG_VIEWBOX = '0 0 45 45'
export const PIECE_STROKE_WIDTH = 1.5

export const PIECE_FILL: Record<Color, string> = {
  [Color.White]: '#fff',
  [Color.Black]: '#000',
}

export const PIECE_STROKE: Record<Color, string> = {
  [Color.White]: '#000',
  [Color.Black]: '#000',
}
