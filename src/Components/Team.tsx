import TopTitle from "./TopTitle";
import TreeItemsTeam from "./Teams/TreeItemsTeam";

function Team({topTitle}:{topTitle:string}){
    return (
        <div style={{maxHeight:"100vh", minHeight:"100vh", 
                backgroundImage:"url('/images/sulakskii.jpg')",
                backgroundSize: "cover",
                backgroundRepeat: "no-repeat",
                backgroundPosition: "center"}}>
        <TopTitle topTitle={topTitle}/>
        <div style={{maxHeight:"90vh", minHeight:"90vh", 
            marginTop:"3vh", marginLeft:"3vw", borderRadius:"30px",
            maxWidth:"79vw", minWidth:"79vw", overflowY:"auto", 
            backgroundColor:"rgba(255, 255, 255, 0.5)",
            fontSize:"14pt",
            }}>
        <TreeItemsTeam />
        </div>
        </div>
    )
}

export default Team;


