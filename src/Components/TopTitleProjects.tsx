import { Calendar2PlusFill } from "react-bootstrap-icons";
import {FC} from 'react'

type DataType = { 
                    id: string;
                    description: string;
                    eventkey: string;
                };

export type ColumnType = {
  topTitle: string;
  setProjects: Function,
  setKey: Function,
  projects: DataType[] | null
};

const TopTitleProjects:FC<ColumnType> = ({topTitle,setProjects, projects,setKey}) => {

    const addNewProject = () => {
        const data_ = new Date().toISOString()
         const foradd = 
            {
                id:`${data_}`,
                description: `${data_} Домашняя страница`,
                eventkey: `${data_} Домашняя страница`,
            }
        const __add = projects
        __add?.push(foradd)
        setProjects(__add)
        // Обновляем ключ и переключаемся на новую вкладку
        setKey(`${data_} Домашняя страница`)
    }

    return (
        <>
        <div className="topTitle whiteText p-3" onClick={addNewProject}>
            <Calendar2PlusFill
            style={{cursor:"pointer"}} />
            <span style={{paddingLeft:"10px"}}>Добавить проект</span>
        {topTitle}</div>
        </>
    )
}

export default TopTitleProjects;