import Modal from 'react-modal';
import { FC } from 'react';
import { Button } from 'react-bootstrap';
// Set the app element for accessibility
Modal.setAppElement('#root');

export type ColumnType = {
  columnId: string;
  deleteColumn: Function,
  setDeleteColumnModalIsOpen: Function,
  DeleteColumnModalIsOpen: boolean
};

var style = {
}

const DeleteColumnModal: FC<ColumnType> = 
({deleteColumn, columnId, setDeleteColumnModalIsOpen, DeleteColumnModalIsOpen}) => {

  return (
    <div>
      <Modal 
        className="modalDelete"
        isOpen={DeleteColumnModalIsOpen}
        onRequestClose={() => setDeleteColumnModalIsOpen(false)}
        contentLabel="DeleteColumnModal"
      >
        <div style={{minHeight:"100px"}}></div>
        <h1>Удалить колонку?</h1>
        <div style={{minHeight:"30px"}}></div>
        <h3>При удалении колонки </h3>
        <h3>будут удалены все задачи в ней</h3>
        <div style={{minHeight:"30px"}}></div>
        <Button 
            style={{
                backgroundColor: "red",
                color: 'white',
                padding: '10px 15px',
                border: 'none',
                borderRadius: '10px',
                fontSize:"18pt",
                cursor: 'pointer',
            }}
            onClick={() => deleteColumn(columnId)}>Удалить</Button>
        <Button style={{
                marginLeft:"10px",
                backgroundColor: "#4f4cec",
                color: 'white',
                padding: '10px 15px',
                border: 'none',
                fontSize:"18pt",
                borderRadius: '10px',

                cursor: 'pointer',}} 
            onClick={() => setDeleteColumnModalIsOpen(false)}>Закрыть</Button>
      </Modal>
    </div>
  );
};

export default DeleteColumnModal;
