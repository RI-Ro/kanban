import { useEffect, useState } from "react";

const EventsList = () => {

    const [events, setEvents] = useState<Node[]>([])

    useEffect(() => {
    setEvents([...initialMockData, ...initialMockData2])
  }, []);


    return (
        <div style={{marginLeft:"40px", marginRight:"40px",
                                    marginBottom:"40px",  marginTop:"40px",
                                    backgroundColor:"rgb(252, 252, 252)",
                                    paddingLeft:"10px", minHeight:"55px",
                                    borderRadius:"25px", paddingBottom:"20px"}}>

            <h1 style={{textAlign:"center", marginBottom:"20px"}}>Лента событий</h1>
            {events &&
            (events.length>0) &&
            
            events.map(node => (
                <div key={node.id} style={{marginLeft:"40px", marginRight:"40px",
                                    marginBottom:"10px",  paddingTop:"10px",
                                    backgroundColor:"rgb(217, 233, 217)",
                                    paddingLeft:"10px", minHeight:"70px",
                                    borderRadius:"25px"}}>
                    {node.tstamp}
                    <strong style={{marginLeft:"10px", marginRight:"10px"}}>{node.type}</strong>
                    {node.name}
                </div>
              ))
            }
        </div>
        
    )
}

export default EventsList;

interface BaseNode {
  id: string;
  type: string;
  name: string;
  tstamp: string;
}

type Node = BaseNode | any

const initialMockData: Node[] = [
  
  // Добавим ещё пользователей для демонстрации (около 100)
  ...Array.from({ length: 10 }, (_, i) => ({
    id: `dev-${i}`,
    type: 'Задачи' as const,
    tstamp: Date().toString(),
    name: `Событие в задачах номер ${i + 1}. Тут может быть любая информация о добавлении, изменении, удалении события и т.д.`,
  })),
];

const initialMockData2: Node[] = [
  
  // Добавим ещё пользователей для демонстрации (около 100)
  ...Array.from({ length: 100 }, (_, i) => ({
    id: `dev-${i}`,
    type: 'Проекты' as const,
    tstamp: Date().toString(),
    name: `Событие в проектах номер ${i + 1}`,
  })),
];