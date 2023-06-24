"use strict";
exports.__esModule = true;
require("./App.css");
function App() {
    return (React.createElement("div", { className: "App", style: {
            display: 'flex',
            flexDirection: 'column',
            height: '100vh',
            width: '100vw',
            justifyContent: 'start',
            alignItems: 'center',
            backgroundColor: '#282c34'
        } },
        React.createElement("h1", { style: {
                color: 'white'
            } }, "Traffic Lights"),
        React.createElement("section", { style: {
                display: 'flex',
                flexDirection: 'row',
                height: '100%',
                width: '100%',
                justifyContent: 'space-between',
                alignItems: 'center'
            } },
            React.createElement("div", { className: 'lights' }),
            React.createElement("div", { className: 'control', style: {
                    display: 'flex',
                    width: '400px',
                    height: '100%',
                    backgroundColor: 'white'
                } }))));
}
exports["default"] = App;
