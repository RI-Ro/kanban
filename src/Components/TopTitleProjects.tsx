import { Calendar2PlusFill } from "react-bootstrap-icons";
import AddProjectsModal from "./AddProjectsModal";
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

const TopTitleProjects:FC<ColumnType> = ({topTitle,setProjects, projects, setKey}) => {

    return (
        <>
        {(projects && projects.length < 10)
        ?
        <AddProjectsModal setProjects={setProjects} topTitle={topTitle} projects={projects} setKey={setKey}/>
        :
        <div className="topTitle whiteText">
        </div>        
        }
        </>
    )
}

export default TopTitleProjects;