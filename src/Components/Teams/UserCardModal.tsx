import Modal from 'react-modal';
import { FC } from 'react';
import { useState } from 'react';
Modal.setAppElement('#root');

export type ColumnType = {
    id : string,
    name : string,
    email? : string,
    phone? : string,
    organization? : string,
    doljnost? : string
};

export type UserType = {
    id : string,
    name : string,
    setIsOpen: Function,
    isOpen: boolean
};

const UserCardModal:FC<UserType> = ({id, name, setIsOpen, isOpen}) => {



    return (

      <Modal 
        className="modalUserCard"
        isOpen={isOpen}
        onRequestClose={() => setIsOpen(false)}
        contentLabel="DeleteColumnModal"
      >
        <div className="card">
  <div className="rounded-top-3" 
  style={{
    backgroundImage: "url(assets/mentor-single.png)",
    backgroundPosition: "center",
    backgroundSize: "cover",
    backgroundRepeat: "no-repeat",
    height: "228px"}}></div>
  <div className="card-body p-md-5">
    <div className="d-flex flex-column gap-5">
      <div className="mt-n8">
        <img src="/assets/default_user.png" alt="mentor 1" 
        className="img-fluid rounded-4 mt-n8" style={{maxHeight:"250px"}}/>
      </div>
      <div className="d-flex flex-column gap-5">
        <div className="d-flex flex-column gap-3">
          <div className="d-flex flex-md-row flex-column justify-content-between gap-2">
            <div>
              <h1 className="mb-0">{name}</h1>
              <div className="d-flex flex-lg-row flex-column gap-2">
                <small className="fw-medium text-gray-800">id={id} Сотрудник ООО "Рога и Копыта"</small>
              </div>
              <div className="d-flex flex-lg-row flex-column gap-2">
                <small className="fw-medium text-success">ООО "Рога и копыта"</small>
              </div>
            </div>
          </div>
          <div className="d-flex flex-md-row flex-column gap-md-4 gap-2">
            <div className="d-flex flex-row gap-2 align-items-center lh-1">
              <span>
                <span className="text-gray-800 fw-bold">EMAIL</span>
              </span>
            </div>
            <div className="d-flex flex-row gap-2 align-items-center lh-1">
              <span>
                <span className="text-gray-800 fw-bold">PHONE</span>
              </span>
            </div>
            <div className="d-flex flex-row gap-2 align-items-center lh-1">
              <span>Написать сообщение</span>
            </div>
          </div>
        </div>
        <div className="d-flex flex-column gap-2">
          <h3 className="mb-0">Участвует в задачах</h3>
          <div className="gap-2 d-flex flex-wrap">
            <a href="#!" className="btn btn-tag btn-sm">Frontend</a>
            <a href="#!" className="btn btn-tag btn-sm">HTML</a>
            <a href="#!" className="btn btn-tag btn-sm">CSS</a>
            <a href="#!" className="btn btn-tag btn-sm">React</a>
            <a href="#!" className="btn btn-tag btn-sm">Javascript</a>
            <a href="#!" className="btn btn-tag btn-sm">Vuejs</a>
            <a href="#!" className="btn btn-tag btn-sm">Next.js</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
      </Modal>

    )
}

export default UserCardModal;