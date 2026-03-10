import { useEffect, useRef, useState } from 'react';
import TopTitleProjects from "./TopTitleProjects"
import Board from "./Boards/Board";
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';

type DataType = { 
                    id: string;
                    description: string;
                    eventkey: string;
                };



function Projects ({topTitle}:{topTitle:string}){
    const scrollContainerRef = useRef<HTMLDivElement>(null);
  
    const [key, setKey] = useState<string | null>('home');

    const [projects, setProjects] = useState<DataType[] | null>(null)

    useEffect(() => {
        setProjects(
        [
            {
                id:"12345678",
                description: 'Домашняя страница',
                eventkey: 'Домашняя страница',
            },
                        {
                id:"23456789",
                description: '2 страница',
                eventkey: '2 страница',
            },
                        {
                id:"34567890",
                description: '3 страница',
                eventkey: '3 страница',
            },
        ]
        )
        // Последняя активная вкладка у пользователя
        setKey('Домашняя страница')     
    }, []); // Empty array means this runs once on mount

    const scrollToLeft= () => {
    if (scrollContainerRef.current) {
        scrollContainerRef.current.scrollTo({
        left:0,
        behavior:"smooth"
        })
        }
    };



    return (
        <>
            <TopTitleProjects topTitle={topTitle} setProjects={setProjects} 
            projects={projects} setKey={setKey}/>
            <div style={{minHeight:"96vh", maxHeight:"96vh", overflowX: "auto", overflowY:"auto",
            backgroundImage:"url('/images/berezovskoe.jpg')",
                 backgroundSize: "cover",
                 backgroundRepeat: "no-repeat",
                 backgroundPosition: "center",
                 backgroundAttachment: "fixed"}}
                 ref={scrollContainerRef}>

    {
        (projects && (projects.length > 0)) ?
      <Tabs 
        style={{backgroundColor:"#535353", fontSize:"16pt"}}
        id="controlled-tabs"
        activeKey={key || ""}
        onSelect={(k) => setKey(k)} // Приведение типа для обработчика
        className="mb-3"
        >
        {
         projects.map((proj)=>(
                <Tab key={proj.id} eventKey={proj.eventkey} title={proj.description}>
                    <Board scrollToLeft={scrollToLeft} project_id={proj.id} />
                </Tab>
            ))
        }
      </Tabs>
        :
      <Tabs
        style={{backgroundColor:"#535353", fontSize:"16pt"}}
        id="controlled-tabs"
        activeKey={key ?? "UNDEFINE__RANDOM_STRING_1234567890_poiuyttrewq"}
        onSelect={(k) => setKey(k)} // Приведение типа для обработчика
        className="mb-3"
      >
      </Tabs>
    }



            </div>
        </>
    )
}

export default Projects;