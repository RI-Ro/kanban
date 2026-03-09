import { useRef } from 'react';
import TopTitle from "./TopTitle"
import Board from "./Boards/Board";

function Projects ({topTitle}:{topTitle:string}){
    const scrollContainerRef = useRef<HTMLDivElement>(null);

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
            <TopTitle topTitle={topTitle}/>
            <div style={{minHeight:"97vh", maxHeight:"97vh", overflowX: "auto", overflowY:"auto",
            backgroundImage:"url('/images/berezovskoe.jpg')",
                 backgroundSize: "cover",
                 backgroundRepeat: "no-repeat",
                 backgroundPosition: "center",
                 backgroundAttachment: "fixed"}}
                 ref={scrollContainerRef}>
            <Board scrollToLeft={scrollToLeft}/>
            </div>
        </>
    )
}

export default Projects;