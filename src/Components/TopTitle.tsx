function TopTitle ({topTitle}:{topTitle:string}) {

    const addNewProject = () => {
        alert("ADD NEW!")
    }

    return (
        <>
        <div className="topTitle whiteText p-3">
        {topTitle}</div>
        </>
    )
}

export default TopTitle;