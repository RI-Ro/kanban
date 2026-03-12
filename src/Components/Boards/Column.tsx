import { FC, useState, useRef, useEffect} from "react";
import { SortableContext, rectSortingStrategy } from "@dnd-kit/sortable";
import { useDroppable } from "@dnd-kit/core";
import Card, { CardType } from "./Card";
import DeleteColumnModal from './DeleteColumnModal'
import { GithubPicker, ColorChangeHandler } from 'react-color';

import { ThreeDotsVertical } from 'react-bootstrap-icons'; // Optional: for icons

export type ColumnType = {
  id: string;
  title: string;
  cards: CardType[];
  borderTop: string;
  background: string;
  deleteColumn: Function
};

const hexTohex= (hex: string): string => {
  // Remove the '#' if present
  // var cleanHex = hex.startsWith('#') ? hex.substring(1) : hex;
  var index = COLORS.indexOf(hex)
  if (index == -1) {
    index = 1
  }
  return BACKGROUNDCOLORS[index]
};

const COLORS = [
  "#83b1cc", "#e9a24f", "#fce258", "#7cae5e", "#49c5bc", "#f75c5c", "#cc8cff", 
  "#b80000", "#008b02", "#5300eb", "#667085", "#eb3737", "#f2732b", "#f5cc00", 
  "#5cdc11", "#08a7a9", "#5089f2", "#e25ef2", "#43f861", "#006b76", "#1273de",
]

const BACKGROUNDCOLORS = [
  "#e3e4ee","#ece7e0","#ecece0","#e0ece0","#e0eceb","#e9dbe6","#e7e0ec",
  "#d3a0a0","#92b893","#bdc4cf","#e3e4ee","#ece0e0","#ece7e0","#ecece0",
  "#e0ece0","#e0eceb","#e0e8ec","#e7e0ec","#d0dfc9", "#bbc6c7", "#b0c8e0"
]


const Column: FC<ColumnType> = ({ id, title, cards, borderTop, background, deleteColumn }) => {
  const [borderTopColor, setborderTopColor] = useState(borderTop);
  const [backgroundColor, setBackgroundColor] = useState(background);
  const [visiblePicker, setVisiblePicker] = useState(false);

  const [DeleteColumnModalIsOpen, setDeleteColumnModalIsOpen] = useState(false);

  const { setNodeRef } = useDroppable({ id: id });

 // Ссылка на контейнер пикера
  const pickerRef = useRef<HTMLDivElement>(null);

  // Обработчик клика вне пикера
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent | TouchEvent) => {
      // Если пикер открыт и клик был вне элемента пикера
      if (
        pickerRef.current &&
        event.target instanceof Node &&
        !pickerRef.current.contains(event.target)
      ) {
        setVisiblePicker(false);
      }
    };
        // Добавляем обработчики только если пикер открыт
    if (visiblePicker) {
      document.addEventListener('mousedown', handleClickOutside);
      document.addEventListener('touchstart', handleClickOutside);
    }

    // Очищаем обработчики при размонтировании или закрытии пикера
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      document.removeEventListener('touchstart', handleClickOutside);
    };
  }, [visiblePicker]);


  const handleChangeColor: ColorChangeHandler = (color) => {
    setBackgroundColor(hexTohex(color.hex));
    setborderTopColor(color.hex);
  };

  return (
    <SortableContext id={id} items={cards} strategy={rectSortingStrategy}>
      <div
        ref={setNodeRef}
        style={{
          marginRight: "30px",
          borderRadius: "30px",
          wordWrap:"break-word",
          borderTop: `solid 10px ${borderTopColor}`,
        }}
      ><div style={{
          minWidth: "250px",
          maxWidth: "250px",
          minHeight: "70px",
          borderRadius: "20px",
          background: `${backgroundColor}`,
        }}>
        <p
          style={{
            padding: "5px 20px",
            textAlign: "left",
            fontWeight: "900",
            color: "#000000"
          }}
        >
          <div className="row">
            <div className="col-10">
              <h4>{title}</h4>
            </div>
            <div className="col-1">
              {visiblePicker ?
              <>
                  <ThreeDotsVertical className="bi me-2" 
                  />
                  <div 
                    ref={pickerRef} 
                    style={{
                      minWidth:"300px", 
                      backgroundColor:"white", 
                      borderRadius:"10px",
                      marginLeft: "-30px",
                      marginTop: "10px",
                      position:"relative",
                      zIndex: "999",
                      }}>
                    <GithubPicker
                    className="TwitterColor"
                    width="190px"
                    colors={COLORS}
                    color={backgroundColor} 
                    onChange={handleChangeColor} />
                    <hr />
                    <div style={{
                      paddingLeft:"10px",
                      color:"red",
                      minWidth:"190px",
                      marginBottom:"10px", 
                      cursor:'pointer',
                      }} 
                      onClick={() => setDeleteColumnModalIsOpen(true)}>Удалить</div>
                      <div style={{minHeight:"10px"}}></div>
                  </div>
                  </>
                  :
                  <ThreeDotsVertical className="bi me-2" 
                    onClick={(event: React.MouseEvent) => {setVisiblePicker(true)}}/>

              }
            </div>
          </div>
        </p>
        {cards.map((card) => (
          <Card key={card.id} id={card.id} title={card.title}></Card>
        ))}
        <div style={{minHeight:"30px"}}></div>
        </div>

<DeleteColumnModal 
deleteColumn={deleteColumn} columnId={id} 
setDeleteColumnModalIsOpen={setDeleteColumnModalIsOpen} 
DeleteColumnModalIsOpen={DeleteColumnModalIsOpen}
/>

      </div>
    </SortableContext>
  );
};

export default Column;
