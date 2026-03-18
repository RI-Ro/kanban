import TopTitle from "./TopTitle";

function Messenger({topTitle}:{topTitle:string}){
    return (
        <>
        <TopTitle topTitle={topTitle}/>
        <div style={{maxHeight:"90vh", minHeight:"90vh", 
            marginTop:"3vh", marginLeft:"3vw", borderRadius:"30px",
            maxWidth:"79vw", minWidth:"79vw", overflowY:"auto", 
            backgroundColor:"rgba(255, 255, 255, 0.5)",
            fontSize:"14pt"}}>
        messenger
        </div>
        </>
    )
}

export default Messenger;