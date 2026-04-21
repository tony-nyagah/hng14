import { Moon, Sun } from 'lucide-react';
import { useTheme } from '../../context/ThemeContext';
import './Sidebar.css';

export default function Sidebar() {
    const { theme, toggleTheme } = useTheme();

    return (
        <aside className="sidebar">
            <div className="sidebar-logo">
                <div className="logo-box">
                    <div className="logo-inner"></div>
                </div>
            </div>

            <div className="sidebar-bottom">
                <button className="theme-toggle" onClick={toggleTheme} aria-label="Toggle theme">
                    {theme === 'light' ? <Moon size={20} color="#7E88C3" /> : <Sun size={20} color="#888EB0" />}
                </button>
                <div className="sidebar-divider"></div>
                <div className="avatar">
                    {/* Using a placeholder avatar for the design */}
                    <img src="https://github.com/shadcn.png" alt="User avatar" />
                </div>
            </div>
        </aside>
    );
}
