// src/components/Sidebar.tsx
import React from 'react';
import {NavLink, Link} from 'react-router-dom'
import { Nav } from 'react-bootstrap';
import { BuildingFill, Calendar2Week, PersonBoundingBox, 
    ChatTextFill, ListColumns, GraphUpArrow, ClipboardDataFill } from 'react-bootstrap-icons'; // Optional: for icons


const Sidebar: React.FC = () => {

const style = {
    color: "white"
}

  return (
    <div className="d-flex flex-column bg-dark p-3" style={{height: '100vh', borderRight:"1px solid #515455" }}>
      <Nav className="nav-pills flex-column mb-auto" >
      <Nav.Item>
            <Nav.Link>
              <NavLink to="/" className='whiteText'
                style={({ isActive }) => ({ paddingLeft: isActive ? '10px' : '0px' })} >
                <ClipboardDataFill className="bi me-2 whiteText" />
                <span className='whiteText'>Проекты</span>
              </NavLink>
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link>
              <NavLink to="/tasks" className='whiteText'
                style={({ isActive }) => ({ paddingLeft: isActive ? '10px' : '0px' })} >
                <Calendar2Week className="bi me-2 whiteText" />
                <span className='whiteText'>Задачи</span>
              </NavLink>
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link>
              <NavLink to="/messenger" className='whiteText'
                style={({ isActive }) => ({ paddingLeft: isActive ? '10px' : '0px' })} >
                <ChatTextFill className="bi me-2 whiteText" />
                <span className='whiteText'>Мессенджер</span>
              </NavLink>
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link>
              <NavLink to="/settings-account" className='whiteText' 
              style={({ isActive }) => ({ paddingLeft: isActive ? '10px' : '0px' })} >
                  <PersonBoundingBox className="bi me-2 whiteText" />
                  <span className='whiteText'>Мой профиль</span>
              </NavLink>
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link>
              <NavLink to="/team" className='whiteText'
                style={({ isActive }) => ({ paddingLeft: isActive ? '10px' : '0px' })} >
                <BuildingFill className="bi me-2 whiteText" />
                <span className='whiteText'>Моя компания</span>
              </NavLink>
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link>
              <NavLink to="/events" className='whiteText'
                style={({ isActive }) => ({ paddingLeft: isActive ? '10px' : '0px' })} >
                <ListColumns className="bi me-2 whiteText" />
                <span className='whiteText'>Лента событий</span>
              </NavLink>
            </Nav.Link>
          </Nav.Item>
          <Nav.Item>
            <Nav.Link>
              <NavLink to="/statistics" className='whiteText'
                style={({ isActive }) => ({ paddingLeft: isActive ? '10px' : '0px' })} >
                <GraphUpArrow className="bi me-2 whiteText" />
                <span className='whiteText'>Отчёты</span>
              </NavLink>
            </Nav.Link>
          </Nav.Item>
      </Nav>
      
    </div>
  );
};

export default Sidebar;
