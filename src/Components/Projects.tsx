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


function Projects ({topTitle}:{topTitle:string}){
    const scrollContainerRef = useRef<HTMLDivElement>(null);
  
    const [key, setKey] = useState<string | null>('home');

//    const [projects, setProjects] = useState<DataType[] | null>(null)
    const [projects, setProjects] = useLocalStorage<DataType[]>(`projects_list`, []);

    useEffect(() => {
        // Последняя активная вкладка у пользователя
        if (projects && (projects.length > 0)) {
            setKey(`${projects[0].id}___${projects[0].eventkey}`)
        } 
    }, []); // Empty array means this runs once on mount

    const scrollToLeft= () => {
    if (scrollContainerRef.current) {
        scrollContainerRef.current.scrollTo({
        left:0,
        behavior:"smooth"
        })
        }
    };

      // Function to handle the deletion
  const handleDeleteProjectByID = (idToDelete:string) => {
    // Create a new array that includes all items EXCEPT the one with the matching id
    if (projects != null) {
        const updatedList = projects.filter((item) => item.id !== idToDelete);
        // Update the state with the new array
        setProjects(updatedList);
        if (projects.length > 0) {
            setKey(`${projects[0].id}___${projects[0].eventkey}`)
        }
    }
    
  };


    return (
        <>
            <div style={{maxHeight:"100vh", minHeight:"100vh", 
                backgroundImage:"url('/images/kislovodsk.jpg')",
                backgroundSize: "cover",
                backgroundRepeat: "no-repeat",
                backgroundPosition: "center"}}>
            <TopTitleProjects topTitle={topTitle} setProjects={setProjects} 
            projects={projects} setKey={setKey}/>

            
            <div style={{maxHeight:"90vh", minHeight:"90vh", 
            marginTop:"3vh", marginLeft:"3vw", borderRadius:"30px",
            maxWidth:"79vw", minWidth:"79vw", overflowY:"auto", 
            backgroundColor:"rgba(255, 255, 255, 0.5)",
            fontSize:"14pt", overflowX: "auto"
                }}
                 ref={scrollContainerRef}>

    {
        (projects && (projects.length > 0)) ?
      <Tabs 
        style={{backgroundColor:"#535353", 
            fontSize:"16pt", width:"300%"}}
        id="controlled-tabs"
        activeKey={key || ""}
        onSelect={(k) => setKey(k)} // Приведение типа для обработчика
        className="mb-3"
        >
        {
         projects.map((proj)=>(
                <Tab key={proj.id} eventKey={`${proj.id}___${proj.eventkey}`} title={proj.description}>
                    <Board scrollToLeft={scrollToLeft} project_id={proj.id} 
                    handleDeleteProjectByID={handleDeleteProjectByID}/>
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

            </div>
        </>
    )
}

export default Projects;

/*
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
*/