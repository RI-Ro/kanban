import { FC, MouseEvent, MouseEventHandler, useState } from "react";
import { CSS } from "@dnd-kit/utilities";
import { useSortable } from "@dnd-kit/sortable";
import { Trash3Fill } from "react-bootstrap-icons";
export type CardType = {
  id: string;
  title: string;
  columnID: string;
  deleteTask: Function;
};

const Card: FC<CardType> = ({ id, title, columnID, deleteTask }) => {
  const { attributes, listeners, setNodeRef, transform } = useSortable({
    id: id
  });

  const onClickDelete = (event: MouseEvent<HTMLDivElement>):void => {
    event.preventDefault()
    deleteTask(id, columnID)
  }

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
  
  const style2 = {
    margin: "10px",
    padding: "5px",
    };

  return (
    <div className="row"  style={{marginLeft:"3px"}}>
    <div className="col-8" ref={setNodeRef} {...attributes} {...listeners} style={style}>
      <div id={id} >
        <p>{title}</p>
      </div>
    </div>
    <div className="col-2" style={style2} onClick={(e) => onClickDelete(e)}>
      <Trash3Fill className="bi me-2" style={{color:"red", cursor: "pointer"}}/></div>
    </div>
    
  );
};

export default Card;
