import { PieceType, Color } from '../lib/chess'

const PIECE_IMAGES: Record<string, string> = {
  'white-K': '/assets/pieces/wk.png',
  'white-Q': '/assets/pieces/wq.png',
  'white-R': '/assets/pieces/wr.png',
  'white-B': '/assets/pieces/wb.png',
  'white-N': '/assets/pieces/wn.png',
  'white-P': '/assets/pieces/wp.png',
  'black-K': '/assets/pieces/bk.png',
  'black-Q': '/assets/pieces/bq.png',
  'black-R': '/assets/pieces/br.png',
  'black-B': '/assets/pieces/bb.png',
  'black-N': '/assets/pieces/bn.png',
  'black-P': '/assets/pieces/bp.png',
}

const PIECE_LETTERS: Record<PieceType, string> = {
  [PieceType.King]: 'K',
  [PieceType.Queen]: 'Q',
  [PieceType.Rook]: 'R',
  [PieceType.Bishop]: 'B',
  [PieceType.Knight]: 'N',
  [PieceType.Pawn]: 'P',
}

export function PieceSVG({ type, color }: { type: PieceType; color: Color }) {
  const src = PIECE_IMAGES[`${color}-${PIECE_LETTERS[type]}`]
  if (!src) return null
  return (
    <img
      src={src}
      alt={`${color} ${type}`}
      draggable={false}
      style={{ width: '100%', height: '100%', objectFit: 'contain' }}
    />
  )
}
