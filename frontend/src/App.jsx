import { useState, createContext } from "react";
import {
  Routes,
  Route,
  Navigate,
  useLocation,
  useNavigate,
} from "react-router";
import { Button, Card, Grid, Header, Menu, Popup } from "semantic-ui-react";
import MainPage from "./pages/Main";
import { ADVICES, PROTOCOL, USERS, HOSTS, PORTS } from "./configuration";

// eslint-disable-next-line react-refresh/only-export-components
export const BackendContext = createContext(null);

function App() {
  const navigate = useNavigate();
  const location = useLocation();

  const backendList = ["Go", "Kotlin", "Java"];

  const backendConfig = [
    {
      users: PROTOCOL + HOSTS[0] + PORTS[0] + USERS,
      advices: PROTOCOL + HOSTS[0] + PORTS[0] + ADVICES,
      active: 0,
    },
    {
      users: PROTOCOL + HOSTS[0] + PORTS[1] + USERS,
      advices: PROTOCOL + HOSTS[0] + PORTS[1] + ADVICES,
      active: 1,
    },
    {
      users: PROTOCOL + HOSTS[0] + PORTS[2] + USERS,
      advices: PROTOCOL + HOSTS[0] + PORTS[2] + ADVICES,
      active: 2,
    },
  ];

  const index = backendList
    .map((i) => i.toUpperCase())
    .indexOf(location.pathname.substring(1).toUpperCase());
  const [backend, setBackend] = useState(backendConfig[Math.max(index, 0)]);

  const changeBackend = (e, { name }) => {
    setBackend(backendConfig[name]);
    navigate(`/${backendList[name].toLowerCase()}`);
  };

  return (
    <Grid columns={3}>
      <Grid.Column width={1} />
      <Grid.Column width={14}>
        <Grid columns={2}>
          <Grid.Row>
            <Card fluid style={{ borderRadius: "0 0 1em 1em" }}>
              <Card.Content>
                <Card.Header>
                  <Grid verticalAlign="middle">
                    <Grid.Column width={4}>
                      <Header size="large">Microservices gRPC</Header>
                    </Grid.Column>
                    <Grid.Column width={8}>
                      <Popup
                        inverted
                        open
                        positionFixed
                        position="right center"
                        trigger={
                          <Menu fluid widths={3} size="huge" compact>
                            {backendList.map((item, index) => (
                              <Menu.Item
                                key={index}
                                style={
                                  index === backend.active
                                    ? { background: "#2185d0", color: "white" }
                                    : {}
                                }
                                active={index === backend.active}
                                onClick={changeBackend}
                                name={index}
                              >
                                {item}
                              </Menu.Item>
                            ))}
                          </Menu>
                        }
                      >
                        Click to change backend!
                      </Popup>
                    </Grid.Column>
                    <Grid.Column width={4} textAlign="right">
                      <Button
                        size="huge"
                        icon="github"
                        color="blue"
                        as="a"
                        target="_blank"
                        href="https://github.com/uid4oe/"
                      />
                      <Button
                        size="huge"
                        icon="linkedin"
                        color="blue"
                        as="a"
                        target="_blank"
                        href="https://www.linkedin.com/in/uid4oe/"
                      />
                    </Grid.Column>
                  </Grid>
                </Card.Header>
              </Card.Content>
            </Card>
          </Grid.Row>
        </Grid>
        <Grid>
          <Grid.Row>
            <Grid.Column>
              <BackendContext.Provider value={backend}>
                <Routes>
                  <Route path="/go" element={<MainPage />} />
                  <Route path="/kotlin" element={<MainPage />} />
                  <Route path="/java" element={<MainPage />} />
                  <Route path="/" element={<Navigate to="/go" replace />} />
                </Routes>
              </BackendContext.Provider>
            </Grid.Column>
          </Grid.Row>
        </Grid>
      </Grid.Column>
      <Grid.Column width={1} />
    </Grid>
  );
}

export default App;
