import {
  closestCorners,
  DndContext,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
  DragOverEvent
} from "@dnd-kit/core";
import { arrayMove, sortableKeyboardCoordinates } from "@dnd-kit/sortable";
import Column, { ColumnType } from "./Column";
import { MouseEventHandler, useEffect, useState } from "react";
import { ArrowLeftCircle, Trash3Fill } from "react-bootstrap-icons";
import DeleteProjectModal from "./DeleteProjectModal";

function useLocalStorage<T>(key: string, initialValue: T): 
[T, (value: T | ((val: T) => T)) => void] {
  const [storedValue, setStoredValue] = useState<T>(() => {
    try {
//      return initialValue;      
      const item = localStorage.getItem(key);
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      return initialValue;
    }
  });

  const setValue = (value: T | ((val: T) => T)) => {
    try {
      const valueToStore = value instanceof Function ? value(storedValue) : value;
      setStoredValue(valueToStore);
      localStorage.setItem(key, JSON.stringify(valueToStore));
    } catch (error) {
      console.error(error);
    }
  };

  return [storedValue, setValue];
}

const Board = ({scrollToLeft, project_id, handleDeleteProjectByID}:{scrollToLeft:MouseEventHandler, project_id:string, handleDeleteProjectByID:Function}) => {
  
  const [DeleteColumnModalIsOpen, setDeleteColumnModalIsOpen] = useState(false);
  
  const data: ColumnType[] = [
    {
      id: "Column1",
      title: "Новые задачи",
      background:"#ecece0",
      borderTop:"#fce258",
      cards: [
      ],
      deleteColumn: Function,
      moveToLeft: Function,
      moveToRight: Function,
      columnPosition: 0,
      isLastPosition: Function,
      changeTitle: Function,
      deleteTask: Function,
    },
    {
      id: "В работе",
      title: "В работе",
      background:"#d0dfc9",
      borderTop:"#43f861",
      cards: [
      ],
      deleteColumn: Function,
      moveToLeft: Function,
      moveToRight: Function,
      columnPosition: 1,
      isLastPosition: Function,
      changeTitle: Function,
      deleteTask: Function,
    },
    {
      id: "Важные",
      title: "Важные",
      background:"#ece0e0",
      borderTop:"#eb3737",
      cards: [
      ],
      deleteColumn: Function,
      moveToLeft: Function,
      moveToRight: Function,
      columnPosition: 2,
      isLastPosition: Function,    
      changeTitle: Function,
      deleteTask: Function,    
    },
    {
      id: "Исполнено",
      title: "Исполнено",
      background:"#e0e8ec",
      borderTop:"#5089f2",
      cards: [
      ],
      deleteColumn: Function,
      moveToLeft: Function,
      moveToRight: Function,
      columnPosition: 3,
      isLastPosition: Function,
      changeTitle: Function,
      deleteTask: Function,    
    }
  ];

//  const [columns, setColumns] = useState<ColumnType[]>(data);

  const [columns, setColumns] = useLocalStorage<ColumnType[]>(`project_id_${project_id}`, data);
  const [addTaskValue, setaddTaskValue] = useState<string>("");
  const [addColumnValue, setaddColumnValue] = useState<string>("");
  const [date, setDate] = useState(Date.now())



  const findColumn = (unique: string | null) => {
    if (!unique) {
      return null;
    }
    if (columns.some((c) => c.id === unique)) {
      return columns.find((c) => c.id === unique) ?? null;
    }
    const id = String(unique);
    const itemWithColumnId = columns.flatMap((c) => {
      const columnId = c.id;
      return c.cards.map((i) => ({ itemId: i.id, columnId: columnId }));
    });
    const columnId = itemWithColumnId.find((i) => i.itemId === id)?.columnId;
    return columns.find((c) => c.id === columnId) ?? null;
  };

  const handleDragOver = (event: DragOverEvent) => {
    const { active, over, delta } = event;
    const activeId = String(active.id);
    const overId = over ? String(over.id) : null;
    const activeColumn = findColumn(activeId);
    const overColumn = findColumn(overId);
    if (!activeColumn || !overColumn || activeColumn === overColumn) {
      return null;
    }
    setColumns((prevState) => {
      const activeItems = activeColumn.cards;
      const overItems = overColumn.cards;
      const activeIndex = activeItems.findIndex((i) => i.id === activeId);
      const overIndex = overItems.findIndex((i) => i.id === overId);
      const newIndex = () => {
        const putOnBelowLastItem =
          overIndex === overItems.length - 1 && delta.y > 0;
        const modifier = putOnBelowLastItem ? 1 : 0;
        return overIndex >= 0 ? overIndex + modifier : overItems.length + 1;
      };
      return prevState.map((c) => {
        if (c.id === activeColumn.id) {
          c.cards = activeItems.filter((i) => i.id !== activeId);
          return c;
        } else if (c.id === overColumn.id) {
          c.cards = [
            ...overItems.slice(0, newIndex()),
            activeItems[activeIndex],
            ...overItems.slice(newIndex(), overItems.length)
          ];
          return c;
        } else {
          return c;
        }
      });
    });
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    const activeId = String(active.id);
    const overId = over ? String(over.id) : null;
    const activeColumn = findColumn(activeId);
    const overColumn = findColumn(overId);
    if (!activeColumn || !overColumn || activeColumn !== overColumn) {
      return null;
    }
    const activeIndex = activeColumn.cards.findIndex((i) => i.id === activeId);
    const overIndex = overColumn.cards.findIndex((i) => i.id === overId);
    if (activeIndex !== overIndex) {
      setColumns((prevState) => {
        return prevState.map((column) => {
          if (column.id === activeColumn.id) {
            column.cards = arrayMove(overColumn.cards, activeIndex, overIndex);
            return column;
          } else {
            return column;
          }
        });
      });
    }
  };

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates
    })
  );

  const AddNewTask = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const c = new FormData(e.currentTarget);
    const addTask = c.get("addTask") as string;
    if (addTask.trim() != "") {
      const updatedData = columns[0].cards;
      const updatedData2 = [...updatedData, {id:Date.now().toString(), title:addTask, columnID:columns[0].id, deleteTask:Function}]
      const updateColumn = columns[0]
      updateColumn.cards = updatedData2
      setColumns(data => [...data.slice(0,0), updateColumn, ...data.slice(1)])
      setaddTaskValue("")
  }
}


const onChangeAddTask = (e: React.ChangeEvent<HTMLInputElement >) => {
   setaddTaskValue(e.target.value)
}

  const AddNewColumn = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const c = new FormData(e.currentTarget);
    const addColumn = c.get("addColumn") as string;
    if (addColumn.trim() != "") {
      setColumns(data => [...data,     {
      id: Date.now().toString(),
      title: addColumn,
      borderTop:"#49c5bc",
      background:"#e0eceb",
      cards: [
      ],
      deleteColumn: Function,
      moveToLeft: Function,
      moveToRight: Function,
      isLastPosition: Function,
      changeTitle: Function,
      columnPosition: 0,
      deleteTask: Function,
      }])
      setaddColumnValue("")
  }
}


const onChangeAddColumn = (e: React.ChangeEvent<HTMLInputElement >) => {
   setaddColumnValue(e.target.value)
}

const deleteColumn = (ColumnIndex:string) => {
  const position = columns.findIndex(item => item.id === ColumnIndex);
  if (position !== -1) {
    setColumns(columns.filter((_, index) => index !== position));
  }
}

const moveToLeft = (ColumnIndex:string) => {
  const position = columns.findIndex(item => item.id === ColumnIndex);
  if (position !== -1) {
    if (position === 0) return; // Некуда двигать, если элемент первый

    const newItems = [...columns];
    // Меняем местами текущий элемент (3) и предыдущий (2)
    [newItems[position - 1], newItems[position]] = [newItems[position], newItems[position - 1]];
    setColumns(newItems);
  }
}

const moveToRight = (ColumnIndex:string) => {
  const position = columns.findIndex(item => item.id === ColumnIndex);
  if (position !== -1) {
    if (position === (columns.length - 1)) return; // Некуда двигать, если элемент последний

    const newItems = [...columns];
    // Меняем местами текущий элемент (3) и предыдущий (2)
    [newItems[position + 1], newItems[position]] = [newItems[position], newItems[position + 1]];
    setColumns(newItems);
  }
}

const isLastPosition = (ColumnPosition:number) => {
  return ColumnPosition === (columns.length-1)
}

const findIndex = (ColumnIndex:string) => {
  return columns.findIndex(item => item.id === ColumnIndex);
}

const deleteProject = () => {
  handleDeleteProjectByID(project_id)
}

const changeTitle = (ColumnIndex:string, newTitle: string) => {
  const position = columns.findIndex(item => item.id === ColumnIndex);
  if (position !== -1) {
    const __columns = columns
    __columns[position].title = newTitle
    setColumns(__columns);
  }
}

const deleteTask = (taskid:string, ColumnIndex:string) => {
  const columnposition = columns.findIndex(item => item.id === ColumnIndex);
  if (columnposition !== -1) {
    const taskposition = columns[columnposition].cards.findIndex(item => item.id === taskid);
    if (taskposition !== -1) {
      const __cards = columns[columnposition].cards.filter((_, index) => index !== taskposition)
      const __columns = columns
      __columns[columnposition].cards = __cards
      setColumns(__columns);
    }
  }
  setDate(Date.now())
}

  return (
    <>
    <DndContext
      sensors={sensors}
      collisionDetection={closestCorners}
      onDragEnd={handleDragEnd}
      onDragOver={handleDragOver}
    >

      <div
        className="Board"
        style={{display: "flex", flexDirection: "row", padding: "20px",            
 }}
      >
      

      
      <div style={{marginRight:"40px",}} >

      <DeleteProjectModal deleteProject={deleteProject}
setDeleteColumnModalIsOpen={setDeleteColumnModalIsOpen} 
DeleteColumnModalIsOpen={DeleteColumnModalIsOpen}
/>

      <div className="deleteProject" onClick={() => setDeleteColumnModalIsOpen(true)}>
        <Trash3Fill size="20px" className="me"/>
        <span style={{marginLeft:"10px"}}>Удалить проект</span></div>

      <form onSubmit={AddNewColumn} style={{paddingTop:"30px"}}>
        <div className="input-box"         
        >
        <label>Создать колонку</label>
        <input name="addColumn" value={addColumnValue} onChange={onChangeAddColumn}/>
        </div>
      </form>
      {
      (columns.length > 0) &&
        <form onSubmit={AddNewTask} style={{paddingTop:"30px"}}>
          <div className="input-box">
          <label>Добавить задачу</label>
          <input name="addTask" value={addTaskValue} onChange={onChangeAddTask}/>
          </div>
        </form>
      }
      </div>

        {columns.map((column) => (
          <Column
            key={column.id}
            id={column.id}
            title={column.title}
            borderTop={column.borderTop}
            background={column.background}
            cards={column.cards}
            deleteColumn={deleteColumn}
            moveToLeft={moveToLeft}
            moveToRight={moveToRight}
            columnPosition={findIndex(column.id)}
            isLastPosition={isLastPosition}
            changeTitle={changeTitle}
            deleteTask={deleteTask}
          ></Column>
        ))}
        
        {
          (columns.length > 4 ) &&
            <div style={{minWidth:"250px"}}
                        onClick={scrollToLeft}
                        >
              <ArrowLeftCircle size="60px" className="me whiteText"/>
            </div>
        }

      </div>
    </DndContext>
    </>
  );
}


export default Board