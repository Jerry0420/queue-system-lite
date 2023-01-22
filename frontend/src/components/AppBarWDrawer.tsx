import React, {useState} from "react"
import PropTypes from "prop-types"
import { Queue } from "../apis/models"
import Box from '@mui/material/Box'
import AppBar from "@mui/material/AppBar"
import Toolbar from "@mui/material/Toolbar"
import IconButton from "@mui/material/IconButton"
import Typography from "@mui/material/Typography"
import Drawer from "@mui/material/Drawer"
import Divider from "@mui/material/Divider"
import ListItem from "@mui/material/ListItem"
import ListItemIcon from "@mui/material/ListItemIcon"
import ListItemText from "@mui/material/ListItemText"
import List from "@mui/material/List"
import MenuIcon from '@mui/icons-material/Menu'
import HomeIcon from '@mui/icons-material/Home'
import HailIcon from '@mui/icons-material/Hail'
import EscalatorWarningIcon from '@mui/icons-material/EscalatorWarning'

const BasicDrawer = (props) => {
    const { setSelectedQueue, queuesInfo, StoreDrawer } = props
    return (
        <div>
        <Toolbar />
        <Divider />
            <ListItem button key={"All"} onClick={() => {setSelectedQueue(null)}}>
                <ListItemIcon>
                <HomeIcon />
                </ListItemIcon>
                <ListItemText primary={"All"} />
            </ListItem>
        <Divider />
        <List>
            {queuesInfo.map((queue: Queue, index: number) => (
            <ListItem button key={queue.id} onClick={() => {setSelectedQueue(queue)}}>
                <ListItemIcon>
                {index % 2 === 0 ? <HailIcon /> : <EscalatorWarningIcon />}
                </ListItemIcon>
                <ListItemText primary={queue.name} />
            </ListItem>
            ))}
        </List>
        {StoreDrawer}
        </div>
    )
}

BasicDrawer.propTypes = {
    setSelectedQueue: PropTypes.func, 
    queuesInfo: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.number,
            name: PropTypes.string,
            customers: PropTypes.arrayOf(
                PropTypes.shape({
                    created_at: PropTypes.string,
                    id: PropTypes.number,
                    name: PropTypes.string,
                    phone: PropTypes.string,
                    state: PropTypes.string,
                })
            )
        })
    ),
    StoreDrawer: PropTypes.node,
}

BasicDrawer.defaultProps = {
    setSelectedQueue: (queue: Queue | null) => {},
    queuesInfo: [],
    StoreDrawer: (<></>),
}

const AppBarWDrawer = (props) => {
    const { storeInfo, children, setSelectedQueue, queuesInfo, StoreDrawer  } = props
    const drawerWidth = 240
    const [mobileOpen, setMobileOpen] = useState(false)
    const handleDrawerToggle = () => {
        setMobileOpen(!mobileOpen)
    }

    return (
        <Box sx={{ display: 'flex' }}>
          <AppBar
            position="fixed"
            sx={{
              width: { sm: `calc(100% - ${drawerWidth}px)` },
              ml: { sm: `${drawerWidth}px` },
            }}
          >
            <Toolbar>
              <IconButton
                color="inherit"
                aria-label="open drawer"
                edge="start"
                onClick={handleDrawerToggle}
                sx={{ mr: 2, display: { sm: 'none' } }}
              >
                <MenuIcon />
              </IconButton>
              <Typography 
                variant="h6" 
                noWrap 
                component="div" 
                style={{cursor: 'pointer'}}
                onClick={() => window.location.reload()}
              >
                {storeInfo.name}
              </Typography>
            </Toolbar>
          </AppBar>
          <Box
            component="nav"
            sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}
            // aria-label="mailbox folders"
          >
            <Drawer
              variant="temporary"
              open={mobileOpen}
              onClose={handleDrawerToggle}
              ModalProps={{
                keepMounted: true,
              }}
              sx={{
                display: { xs: 'block', sm: 'none' },
                '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
              }}
            >
              <BasicDrawer
                setSelectedQueue={setSelectedQueue}
                queuesInfo={queuesInfo}
                StoreDrawer={StoreDrawer}
              />
            </Drawer>
            <Drawer
              variant="permanent"
              sx={{
                display: { xs: 'none', sm: 'block' },
                '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
              }}
              open
            >
              <BasicDrawer 
                setSelectedQueue={setSelectedQueue}
                queuesInfo={queuesInfo}
                StoreDrawer={StoreDrawer}
              />
            </Drawer>
          </Box>
          <Box
            component="main"
            sx={{ flexGrow: 1, p: 3, width: { sm: `calc(100% - ${drawerWidth}px)` } }}
          >
            <Toolbar />
            {children}
          </Box>
        </Box>
      )
}

AppBarWDrawer.propTypes = {
    storeInfo: PropTypes.shape({
        name: PropTypes.string
    }),
    children: PropTypes.node,
    setSelectedQueue: PropTypes.func, 
    queuesInfo: PropTypes.arrayOf(
        PropTypes.shape({
            id: PropTypes.number,
            name: PropTypes.string,
            customers: PropTypes.arrayOf(
                PropTypes.shape({
                    created_at: PropTypes.string,
                    id: PropTypes.number,
                    name: PropTypes.string,
                    phone: PropTypes.string,
                    state: PropTypes.string,
                })
            )
        })
    ),
    StoreDrawer: PropTypes.node,
}

AppBarWDrawer.defaultProps = {
    storeInfo: {"name": ""},
    children: (<></>),
    setSelectedQueue: (queue: Queue | null) => {},
    queuesInfo: [],
    StoreDrawer: (<></>),
}

export {
    AppBarWDrawer
}