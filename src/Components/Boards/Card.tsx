import { FC } from "react";
import { CSS } from "@dnd-kit/utilities";
import { useSortable } from "@dnd-kit/sortable";
export type CardType = {
  id: string;
  title: string;
};

const Card: FC<CardType> = ({ id, title }) => {
  const { attributes, listeners, setNodeRef, transform } = useSortable({
    id: id
  });

  const style = {
    margin: "10px",
    opacity: 1,
    color: "#333",
    background: "white",
    padding: "5px",
    borderRadius: "10px",
    fontSize:"14pt",
    transform: CSS.Transform.toString(transform)
  };

  return (
    <div ref={setNodeRef} {...attributes} {...listeners} style={style}>
      <div id={id}>
        <p>{title}</p>
      </div>
    </div>
  );
};

export default Card;
