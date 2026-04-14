import { useState, useMemo, useCallback, useRef, useEffect } from 'react'
import {
  parseFEN,
  parseFENMeta,
  getPseudoLegalMoves,
  toAlgebraic,
  Pos,
  BOARD_SIZE,
  isInCheck,
  findKing,
} from '../../lib/chess'
import { PieceType, Color } from '../../lib/chess'
import { PieceSVG } from '../ChessPieces'
import * as S from './styles'
import { DRAG_THRESHOLD_PX, DRAG_OVERLAY_SIZE } from '../../lib/config/drag'
import { GameStatus } from '../../lib/game'

interface ChessBoardProps {
  fen: string
  turn: Color
  status: string
  onMove: (from: string, to: string) => void
  flipped?: boolean
}

interface DragState {
  fromRow: number
  fromCol: number
  pieceType: PieceType
  pieceColor: Color
  startX: number
  startY: number
  x: number
  y: number
}

const ROWS_WHITE = [0, 1, 2, 3, 4, 5, 6, 7]
const ROWS_BLACK = [7, 6, 5, 4, 3, 2, 1, 0]

export function ChessBoard({ fen, turn, status, onMove, flipped = false }: ChessBoardProps) {
  const board = useMemo(() => parseFEN(fen), [fen])
  const fenMeta = useMemo(() => parseFENMeta(fen), [fen])
  const inCheck = useMemo(() => isInCheck(board, turn), [board, turn])
  const kingPos = useMemo(() => (inCheck ? findKing(board, turn) : null), [inCheck, board, turn])

  const rows = flipped ? ROWS_BLACK : ROWS_WHITE
  const cols = flipped ? ROWS_BLACK : ROWS_WHITE

  const [selected, setSelected] = useState<Pos | null>(null)
  const [validMoves, setValidMoves] = useState<Pos[]>([])
  const [captures, setCaptures] = useState<Set<string>>(new Set())
  const [drag, setDrag] = useState<DragState | null>(null)
  const [dragOver, setDragOver] = useState<Pos | null>(null)

  const dragRef = useRef<DragState | null>(null)
  const validMovesRef = useRef(validMoves)
  const resolveRef = useRef<(x: number, y: number) => Pos | null>(() => null)
  const didDrag = useRef(false)

  validMovesRef.current = validMoves

  const coordKey = (r: number, c: number) => `${r},${c}`

  const selectPiece = useCallback(
    (row: number, col: number) => {
      setSelected({ row, col })
      const moves = getPseudoLegalMoves(board, row, col, turn, fenMeta)
      setValidMoves(moves)
      const caps = new Set<string>()
      for (const m of moves) {
        if (board[m.row][m.col]) caps.add(coordKey(m.row, m.col))
      }
      setCaptures(caps)
    },
    [board, turn, fenMeta],
  )

  const clearSelection = useCallback(() => {
    setSelected(null)
    setValidMoves([])
    setCaptures(new Set())
  }, [])

  const executeMove = useCallback(
    (fromRow: number, fromCol: number, toRow: number, toCol: number) => {
      onMove(toAlgebraic(fromRow, fromCol), toAlgebraic(toRow, toCol))
      clearSelection()
    },
    [onMove, clearSelection],
  )

  const resolveSquare = useCallback(
    (clientX: number, clientY: number): Pos | null => {
      const el = document.querySelector('[data-board]') as HTMLElement | null
      if (!el) return null
      const rect = el.getBoundingClientRect()
      const size = rect.width / BOARD_SIZE
      const vx = Math.floor((clientX - rect.left) / size)
      const vy = Math.floor((clientY - rect.top) / size)
      if (vx < 0 || vx >= BOARD_SIZE || vy < 0 || vy >= BOARD_SIZE) return null
      return flipped ? { row: BOARD_SIZE - 1 - vy, col: BOARD_SIZE - 1 - vx } : { row: vy, col: vx }
    },
    [flipped],
  )

  useEffect(() => {
    resolveRef.current = resolveSquare
  }, [resolveSquare])

  useEffect(() => {
    if (!drag) return

    const onMove = (e: PointerEvent) => {
      dragRef.current = { ...dragRef.current!, x: e.clientX, y: e.clientY }
      setDrag(prev => (prev ? { ...prev, x: e.clientX, y: e.clientY } : null))
      setDragOver(resolveRef.current(e.clientX, e.clientY))
    }

    const onUp = (e: PointerEvent) => {
      const d = dragRef.current!
      const dx = e.clientX - d.startX
      const dy = e.clientY - d.startY
      if (dx * dx + dy * dy > DRAG_THRESHOLD_PX * DRAG_THRESHOLD_PX) didDrag.current = true

      const target = resolveRef.current(e.clientX, e.clientY)
      if (target && didDrag.current) {
        const isValid = validMovesRef.current.some(
          m => m.row === target.row && m.col === target.col,
        )
        if (isValid) {
          executeMove(d.fromRow, d.fromCol, target.row, target.col)
        }
      }
      dragRef.current = null
      setDrag(null)
      setDragOver(null)
    }

    document.addEventListener('pointermove', onMove)
    document.addEventListener('pointerup', onUp)
    return () => {
      document.removeEventListener('pointermove', onMove)
      document.removeEventListener('pointerup', onUp)
    }
  }, [drag, executeMove])

  const onSquareClick = (row: number, col: number) => {
    if (status !== GameStatus.Active) return
    if (didDrag.current) {
      didDrag.current = false
      return
    }
    const piece = board[row][col]

    if (selected) {
      const isValid = validMoves.some(m => m.row === row && m.col === col)
      if (isValid) {
        executeMove(selected.row, selected.col, row, col)
        return
      }
      if (piece && piece.color === turn) {
        selectPiece(row, col)
        return
      }
      clearSelection()
      return
    }

    if (piece && piece.color === turn) selectPiece(row, col)
  }

  const onPointerDown = (e: React.PointerEvent, row: number, col: number) => {
    if (status !== GameStatus.Active) return
    const piece = board[row][col]
    if (!piece || piece.color !== turn) return

    selectPiece(row, col)
    const info: DragState = {
      fromRow: row,
      fromCol: col,
      pieceType: piece.type,
      pieceColor: piece.color,
      startX: e.clientX,
      startY: e.clientY,
      x: e.clientX,
      y: e.clientY,
    }
    dragRef.current = info
    setDrag(info)
    ;(e.target as HTMLElement).setPointerCapture(e.pointerId)
  }

  const isMoveTarget = (r: number, c: number) => validMoves.some(m => m.row === r && m.col === c)
  const isCaptureTarget = (r: number, c: number) => captures.has(coordKey(r, c))
  const showDragOverlay =
    drag !== null && Math.hypot(drag.x - drag.startX, drag.y - drag.startY) > DRAG_THRESHOLD_PX

  return (
    <>
      <S.Board data-board $dragging={showDragOverlay}>
        {rows.map(ri =>
          cols.map(ci => {
            const piece = board[ri][ci]
            const light = (ri + ci) % 2 === 0
            const sel = selected?.row === ri && selected?.col === ci
            const dragSrc = showDragOverlay && drag.fromRow === ri && drag.fromCol === ci
            const dragOv = showDragOverlay && dragOver?.row === ri && dragOver?.col === ci
            const inCheckSquare = inCheck && kingPos?.row === ri && kingPos?.col === ci

            const SquareComponent = inCheckSquare ? S.KingCheckSquare : S.Square

            return (
              <SquareComponent
                key={`${ri}-${ci}`}
                $isLight={light}
                $selected={sel}
                $dragSource={dragSrc}
                $dragOver={dragOv}
                $kingCheck={inCheckSquare}
                onClick={() => onSquareClick(ri, ci)}
                onPointerDown={ev => onPointerDown(ev, ri, ci)}
              >
                {piece && <PieceSVG type={piece.type} color={piece.color} />}
                {isMoveTarget(ri, ci) && !isCaptureTarget(ri, ci) && <S.MoveIndicator />}
                {isCaptureTarget(ri, ci) && <S.CaptureIndicator />}
              </SquareComponent>
            )
          }),
        )}
      </S.Board>

      {showDragOverlay && drag && (
        <S.DraggedPiece style={{ left: drag.x - DRAG_OVERLAY_SIZE / 2, top: drag.y - DRAG_OVERLAY_SIZE / 2 }}>
          <PieceSVG type={drag.pieceType} color={drag.pieceColor} />
        </S.DraggedPiece>
      )}
    </>
  )
}
