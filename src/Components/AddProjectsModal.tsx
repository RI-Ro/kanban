import Modal from 'react-modal';
import { FC, ChangeEvent, KeyboardEvent } from 'react';
import { Button } from 'react-bootstrap';
import {useState, useRef, useEffect} from 'react';
import { Calendar2PlusFill } from 'react-bootstrap-icons';
// Set the app element for accessibility
Modal.setAppElement('#root');

export type ColumnType = {
  topTitle: string;
  setProjects: Function,
  projects: DataType[] | null,
  setKey: Function,
};

type DataType = { 
                    id: string;
                    description: string;
                    eventkey: string;
                };

var style = {
}

const AddProjectsModal: FC<ColumnType> = 
({setProjects, topTitle, projects, setKey}) => {

    const [isOpen, setIsOpen] = useState(false)
    const [addText, setAddText] = useState("")
    const inputRef = useRef<HTMLInputElement>(null);
  
    useEffect(() => {
        if (isOpen) {
        // Таймаут или requestAnimationFrame иногда нужны, 
        // если модалка анимирована (например, через CSS transitions)
        const timer = setTimeout(() => inputRef.current?.focus(), 100);
        return () => clearTimeout(timer);
        }
    }, [isOpen]);

    const addProjects = () => {
        addNewProject()
    }

        const addNewProject = () => {
        const data_ = new Date().toISOString()
         const foradd = 
            {
                id:`${data_}`,
                description: addText.trim(),
                eventkey: addText.trim(),
            }
        const __add = projects
        __add?.push(foradd)
        setProjects(__add)
        // Обновляем ключ и переключаемся на новую вкладку
        setKey(addText.trim())
        setAddText("")
        setIsOpen(false)
    }

    const setIsOpenModal = () => {
        setIsOpen(true)
    }

    const onchangeInput = (e: ChangeEvent<HTMLInputElement>) => {
        e.preventDefault()
        setAddText(e.target.value)
    }
  
  // Handle key presses, specifically the 'Enter' key
  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      // Prevent default form submission behavior if inside a form
      e.preventDefault();
      if (addText.trim().length > 0) {
        addProjects()
      }
      // You can add your submit logic or function call here
      // e.g., handleSubmit(value);
    }
  };
  return (

    <>
        <div className="topTitle whiteText" onClick={setIsOpenModal}
        style={{paddingLeft:"20px", paddingTop:"5px"}}>
            <Calendar2PlusFill
            style={{cursor:"pointer"}} />
            <span style={{paddingLeft:"10px"}}>Добавить проект</span>
        {topTitle}
        </div>

    <div>
      <Modal 
        className="modalAdd"
        isOpen={isOpen}
        onRequestClose={() => setIsOpen(false)}
        contentLabel="DeleteColumnModal"
      >
        <div style={{minHeight:"100px"}}></div>
        <form>
            <input type="text" name='addInput' value={addText} onChange={onchangeInput} 
            onKeyDown={handleKeyDown} ref={inputRef}  maxLength={25}/>
        </form>
        <div style={{minHeight:"30px"}}></div>
        <h3>Введите описание нового проекта</h3>
        <p>* описание проекта не более 25 символов c учетом пробелов<br />
        ** Можно добавить не более 10 проектов</p>
        <div style={{minHeight:"30px"}}></div>
        <div className='row'>
            <div className='col-6'></div>
        { (addText.trim().length > 0) ?
        <div className='col-3'>
            <Button 
            style={{
                backgroundColor: "green",
                color: 'white',
                padding: '10px 15px',
                border: 'none',
                borderRadius: '10px',
                fontSize:"18pt",
                cursor: 'pointer',
            }}
            onClick={() => addProjects()}>Добавить</Button>
        </div>
            :
        <div className='col-3'></div>
        }
        <div className='col-3'>
        <Button style={{
                marginLeft:"10px",
                backgroundColor: "#4f4cec",
                color: 'white',
                padding: '10px 15px',
                border: 'none',
                fontSize:"18pt",
                borderRadius: '10px',
                cursor: 'pointer',}} 
            onClick={() => setIsOpen(false)}>Закрыть</Button>
        </div>
        </div>
      </Modal>

    </div>
        </>

  );
};

export default AddProjectsModal;


