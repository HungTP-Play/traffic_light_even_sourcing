import React from 'react';
import shortUUID from 'short-uuid';
import { post } from 'superagent';

import greenImg from './traffic_light_images/green.png';
import redImg from './traffic_light_images/red.png';
import yellowImg from './traffic_light_images/yellow.png';

const URL = 'http://localhost:1111/emit';

function getEventId(): string {
    const uuid = shortUUID.generate();
    return `event-${uuid}`
}

function registerLight(lightId: string, location: string) {
    post(URL).send({
        event_id: getEventId(),
        event_name: 'registration_event',
        event_data: {
            light_id: lightId,
            location: location,
        }
    }).then((res: any) => {
        console.log(res);
    });
}

function stateChange(lightId: string, location: string, prevColor: string, nextColor: string) {
    post(URL).send({
        event_id: getEventId(),
        event_name: 'state_change_event',
        event_data: {
            light_id: lightId,
            location: location,
            from_state: prevColor,
            to_state: nextColor,
        }
    }).then((res: any) => {
        console.log(res);
    });
}
export function Light(props: {
    times: [number, number, number], index: number,
    latLong: [number, number],
    lightId: string,
}) {
    const [lat, long] = props.latLong;
    const location = `${lat}::${long}`;
    const [red, yellow, green] = props.times;
    const [color, setColor] = React.useState('RED');

    const colorImageLink = React.useMemo(() => {
        switch (color) {
            case 'RED':
                return redImg;
            case 'YELLOW':
                return yellowImg;
            case 'GREEN':
                return greenImg;
            default:
                return redImg;
        }
    }, [color]);

    // Create timer (interval) for each color
    // RED, YELLOW, GREEN, YELLOW, RED,...
    React.useEffect(() => {
        const firstDelay = 1 * 1000;
        let interval: NodeJS.Timeout;
        const firstDelayTimeout = setTimeout(() => {
            const intervalTime = (red + yellow * 2 + green) * 1000;
            interval = setInterval(() => {
                stateChange(props.lightId, location, 'YELLOW', 'RED');
                setColor('RED');

            }, intervalTime);
        }
            , firstDelay);
        return () => {
            clearTimeout(firstDelayTimeout);
            clearInterval(interval);
        }
    }, [red, green, yellow, props.lightId, location]);

    // YELLOW 1
    React.useEffect(() => {
        const firstDelay = red * 1000 + 1 * 1000;
        let yellowInterval: NodeJS.Timeout;
        const firstDelayTimeout = setTimeout(() => {
            const intervalTime = (green + yellow + red) * 1000;
            yellowInterval = setInterval(() => {
                stateChange(props.lightId, location, 'RED', 'YELLOW');
                setColor('YELLOW');
            }, intervalTime);
        }, firstDelay);
        return () => {
            clearTimeout(firstDelayTimeout);
            clearInterval(yellowInterval);
        }
    }, [red, yellow, green, props.lightId, location]);

    // GREEN
    React.useEffect(() => {
        const firstDelay = red * 1000 + yellow * 1000 + 1 * 1000;
        let greenInterval: NodeJS.Timeout;
        const delayTimeout = setTimeout(() => {
            const intervalTime = (red + yellow * 2) * 1000;
            greenInterval = setInterval(() => {
                stateChange(props.lightId, location, 'YELLOW', 'GREEN');
                setColor('GREEN');
            }, intervalTime);
        }, firstDelay);
        return () => {
            clearTimeout(delayTimeout);
            clearInterval(greenInterval);
        }
    }, [red, green, yellow, props.lightId, location]);

    // YELLOW 2
    React.useEffect(() => {
        const firstDelay = red * 1000 + yellow * 1000 + green * 1000 + 1 * 1000;
        let yellowInterval: NodeJS.Timeout;
        const delayTimeout = setTimeout(() => {
            const intervalTime = (red + green + yellow * 2) * 1000;
            yellowInterval = setInterval(() => {
                stateChange(props.lightId, location, 'GREEN', 'YELLOW');
                setColor('YELLOW');
            }, intervalTime);
        }, firstDelay);
        return () => {
            clearTimeout(delayTimeout);
            clearInterval(yellowInterval);
        }
    }, [red, green, yellow, props.lightId, location]);

    // Call the registration endpoint only once
    React.useEffect(() => {
        registerLight(props.lightId, location);
    }, [props.lightId, location]);

    return <div style={{
        display: 'flex',
        flexDirection: 'column',
        height: '150px',
        width: '50px',
    }} key={props.index}>
        <img src={colorImageLink} alt={`traffic_light_${props.lightId}_${color}`} style={{
            height: '100%',
            width: '100%',
            objectFit: 'cover',
        }} />
    </div>
}