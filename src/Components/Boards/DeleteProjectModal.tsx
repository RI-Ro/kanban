import Modal from 'react-modal';
import { FC } from 'react';
import { Button } from 'react-bootstrap';
// Set the app element for accessibility
Modal.setAppElement('#root');

export type ColumnType = {
  deleteProject: Function,
  setDeleteColumnModalIsOpen: Function,
  DeleteColumnModalIsOpen: boolean
};

var style = {
}

const DeleteProjectModal: FC<ColumnType> = 
({deleteProject, setDeleteColumnModalIsOpen, DeleteColumnModalIsOpen}) => {

  return (
    <div>
      <Modal 
        className="modalDelete"
        isOpen={DeleteColumnModalIsOpen}
        onRequestClose={() => setDeleteColumnModalIsOpen(false)}
        contentLabel="DeleteColumnModal"
      >
        <div style={{minHeight:"100px"}}></div>
        <h1>Удалить проект?</h1>
        <div style={{minHeight:"30px"}}></div>
        <h3>При удалении проекта</h3>
        <h3>будут удалены колонки и все задачи в них!</h3>
        <div style={{minHeight:"30px"}}></div>
                <div className='row'>
            <div className='col-7'></div>
        <div className='col-2'>
        <Button className='customBtnDelete'
            onClick={() => deleteProject()}>Удалить</Button>
        </div>
        <div className='col-2'>
        <Button className='customBtnClose'
            onClick={() => setDeleteColumnModalIsOpen(false)}>Закрыть</Button>
        </div>
        </div>


      </Modal>
    </div>
  );
};

export default DeleteProjectModal;
