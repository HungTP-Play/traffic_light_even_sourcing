import * as React from 'react';
import shortUUID from 'short-uuid';
import './App.css';
import { Light } from './Light';
function randomLatLong(): [number, number] {
  const lat = randomIntBetween(-90, 90);  // Latitude between -90 to 90 degrees
  const long = randomIntBetween(-180, 180); // Longitude between -180 to 180 degrees

  return [lat, long];
}

function randomIntBetween(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}


function getUUID(): string {
  const uuid = shortUUID.generate();
  return `light-${uuid}`
}

function App() {
  const [lights, setLights] = React.useState<{
    lightId: string
    times: [number, number, number]
    latLong: [number, number]
  }[]>([]) // [{times:[red, yellow, green]}]

  const greeRef = React.useRef<HTMLInputElement>(null);
  const yellowRef = React.useRef<HTMLInputElement>(null);
  const redRef = React.useRef<HTMLInputElement>(null);

  const onCreateLight = (event: any) => {
    event.preventDefault()
    const latLong = randomLatLong()
    const lightId = getUUID()

    const redInt = parseInt(redRef.current?.value ?? '0')
    const yellowInt = parseInt(yellowRef.current?.value ?? '0')
    const greenInt = parseInt(greeRef.current?.value ?? '0')
    const newLight: {
      lightId: string
      times: [number, number, number]
      latLong: [number, number]
    } = {
      lightId: lightId,
      times: [redInt, yellowInt, greenInt],
      latLong: [latLong[0], latLong[1]]
    }
    setLights([...lights, newLight])
    event.target.reset()
  }

  return (
    <div className="App" style={{
      display: 'flex',
      flexDirection: 'column',
      height: '100vh',
      width: '100vw',
      justifyContent: 'start',
      alignItems: 'center',
      backgroundColor: '#282c34',
    }}>
      <h1 style={{
        color: 'white',
      }}>Traffic Lights</h1>

      <section style={{
        display: 'flex',
        flexDirection: 'row',
        height: '100%',
        width: '100%',
        justifyContent: 'space-between',
        alignItems: 'center',
      }}>

        <div className='lights' style={{
          display: 'flex',
          flexDirection: 'row',
          justifyContent: 'flex-start',
          alignItems: 'flex-start',
          height: '100%',
          width: 'calc(100% - 400px)',
        }}>
          {lights.map((light, index) => {
            return <Light latLong={
              light.latLong
            } lightId={light.lightId} times={light.times} index={index} key={`${getUUID()}`} />
          })}
        </div>
        <div className='control' style={{
          display: 'flex',
          width: '400px',
          height: '100%',
          backgroundColor: '#282c34',
        }}>
          <form style={{
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'flex-start',
            alignItems: 'flex-start',
            height: '100%',
            width: '100%',
            padding: '0px 5px',
          }} onSubmit={onCreateLight}>

            <label style={{
              fontSize: '20px',
              display: 'inline-block',
              color: 'white',
            }}>Light Duration</label>
            <span style={{ height: '10px' }}></span>
            <label style={{
              color: 'red',
            }}>Red</label>
            <input type='number' placeholder='Red' style={{
              width: '90%',
              height: '50px',
              borderRadius: '5px',
              padding: '0px 10px',
            }} name='red_time' ref={redRef} />

            <span style={{ height: '10px' }}></span>
            <label style={{
              color: 'yellow',
            }}>Yellow</label>
            <input type='number' placeholder='Yellow' style={{
              width: '90%',
              height: '50px',
              borderRadius: '5px',
              padding: '0px 10px',
            }} name='yellow_time' ref={yellowRef}/>

            <span style={{ height: '10px' }}></span>
            <label style={{
              color: 'green',
            }}>Green</label>
            <input type='number' placeholder='Green' style={{
              width: '90%',
              height: '50px',
              borderRadius: '5px',
              padding: '0px 10px',
            }} name='green_time' ref={greeRef} />

            <button style={{
              alignSelf: 'center',
              marginTop: '20px',
              backgroundColor: 'green',
              color: 'white',
              fontSize: '20px',
              padding: '10px 20px',
            }} type='submit'>
              Create light
            </button>
          </form>
        </div>
      </section>
    </div>
  );
}

export default App;
