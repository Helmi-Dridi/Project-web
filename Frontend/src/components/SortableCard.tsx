import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import type { JSX } from "react";

type SortableCardProps = {
  id: string;
  children: (dragHandleProps: {
    attributes?: React.HTMLAttributes<HTMLElement>;
    listeners?: React.HTMLAttributes<HTMLElement>;
  }) => JSX.Element;
};

export function SortableCard({ id, children }: SortableCardProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
  } = useSortable({ id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div ref={setNodeRef} style={style}>
      {children({
        attributes,
        listeners,
      })}
    </div>
  );
}
