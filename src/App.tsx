
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Sidebar from "./Components/Sidebar";
import Blank from './Components/Blank';
import 'bootstrap/dist/css/bootstrap.min.css';
import './styles.css'
import Projects from "./Components/Projects";
import SettingsAccount from "./Components/SettingsAccount";
import Team from "./Components/Team";
import Tasks from "./Components/Tasks";
import Messenger from "./Components/Messenger";
import Events from "./Components/Events";
import Statistics from "./Components/Statistic";

const App: React.FC = () => {
  return (
    <div >
        <BrowserRouter>
        <div className="block1">
            <Sidebar />
        </div>
        <div className="block2" style={{backgroundImage:"url('/images/background.jpg')", backgroundSize: "cover"}}>
        <Routes>
            <Route path="/settings-account" element={<SettingsAccount topTitle=""/>} />
            <Route path="/team" element={<Team topTitle=""/>} />
            <Route path="/" element={<Projects topTitle=""/>} />
            <Route path="/tasks" element={<Tasks topTitle=""/>} />
            <Route path="/messenger" element={<Messenger topTitle=""/>} />
            <Route path="/events" element={<Events topTitle=""/>} />
            <Route path="/statistics" element={<Statistics topTitle=""/>} />
        </Routes>
        </div>
        <div className="block3"></div>
        </BrowserRouter>
    </div>
  );
}

export default App;


