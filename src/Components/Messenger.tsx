import TopTitle from "./TopTitle";

function Messenger({topTitle}:{topTitle:string}){
    return (
        <>
        <TopTitle topTitle={topTitle}/>
        Messenger
        </>
    )
}

export default Messenger;