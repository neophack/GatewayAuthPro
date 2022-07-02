import './App.css';

import {useState} from 'react';
import {withRouter} from 'react-router-dom'
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Container from '@mui/material/Container';
import Paper from '@mui/material/Paper';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import Backdrop from '@mui/material/Backdrop';
import CircularProgress from '@mui/material/CircularProgress';
import {useSnackbar} from 'notistack';
import {Login} from "./utils/request";
import {getUrlParams, isBlank} from './utils/utils'
import logo from './image/logo.png'
import bgimage from './image/1920_middle.jpg'


function App(props) {

    const [account, setAccount] = useState("");
    const [password, setPassword] = useState("");
    const [backdropOpen, setBackdropOpen] = useState(false);
    const {enqueueSnackbar} = useSnackbar();

    const snackbarStype = (variant) => {
        return {
            variant: variant,
            anchorOrigin: {
                vertical: 'top', horizontal: 'center'
            }
        }
    }

    const handleSubmit = () => {
        var url = getUrlParams("url", props.location.search);
        if (isBlank(url)) {
            url = "/"
        }
        if (isBlank(account)) {
            enqueueSnackbar("Account cannot be empty", snackbarStype('error'));
            return
        }
        if (isBlank(password)) {
            enqueueSnackbar("Password cannot be empty", snackbarStype('error'));
            return
        }
        setBackdropOpen(true);
        Login("api/loginxx", account, password)
            .then(res => {
                if (res.code === 200) {
                    enqueueSnackbar(res.msg, snackbarStype('success'));
                    window.location.href = url
                } else {
                    enqueueSnackbar(res.msg, snackbarStype('error'));
                }
            }).finally(() => {
            setBackdropOpen(false);
        })
    };

    const onKeyPress = (e) => {
        if (e.key === "Enter") {
            handleSubmit()
        }
    };

    const handleSubmit2 = () => {
        var url = getUrlParams("url", props.location.search);
        window.location.href = "/logout?url=" + url
    }

    const html = () => {
        var type = getUrlParams("type", props.location.search);
        if (type === 'nopermission') {
            return (
                <Box
                    sx={{
                        display: 'flex',
                        flexDirection: 'column',
                        minHeight: '100vh',
                    }}

                >
                    <CssBaseline/>
                    <Container component="main" sx={{mt: 8, mb: 2}} maxWidth="sm">
                        <Typography variant="h4" component="h1" gutterBottom>
                            Accounts no permissions,<br/>Please log in again!<br/>
                            账号无权限，<br/>请重新登录!<br/>
                        </Typography>
                        <Typography variant="body1">
                            <Button
                                variant="contained"
                                color="error"
                                sx={{mt: 3, mb: 2}}
                                onClick={handleSubmit2}
                                size="large"
                            >
                                Re-register / 重新登录
                            </Button>
                        </Typography>
                    </Container>
                </Box>
            )
        } else {
            return (
                <Grid container component="main" sx={{height: '100vh'}}>
                    <CssBaseline/>
                    <Grid
                        item
                        xs={false}
                        sm={4}
                        md={7}
                        sx={{
                            backgroundImage: 'url('+bgimage+')',
                            backgroundRepeat: 'no-repeat',
                            backgroundColor: (t) =>
                                t.palette.mode === 'light' ? t.palette.grey[50] : t.palette.grey[900],
                            backgroundSize: 'cover',
                            backgroundPosition: 'center',
                        }}
                    />
                    <Grid item xs={12} sm={8} md={5} component={Paper} elevation={6} square>
                        <Box
                            sx={{
                                my: 8,
                                mx: 4,
                                display: 'flex',
                                flexDirection: 'column',
                                alignItems: 'center',
                            }}
                        >
                            <Avatar sx={{m: 1, bgcolor: '#fff'}}>
                                <img
                                    src={logo}
                                    alt={"logo"}
                                    width={'100%'}
                                    loading="lazy"
                                />
                            </Avatar>
                            <Typography component="h1" variant="h5">
                                Login
                            </Typography>
                            <Box component="form" noValidate sx={{mt: 1}}>
                                <TextField
                                    margin="normal"
                                    required
                                    fullWidth
                                    id="account"
                                    label="account"
                                    name="account"
                                    autoFocus
                                    value={account}
                                    onChange={(event) => {
                                        setAccount(event.target.value)
                                    }}
                                />
                                <TextField
                                    margin="normal"
                                    required
                                    fullWidth
                                    name="password"
                                    label="password"
                                    type="password"
                                    id="password"
                                    autoComplete="current-password"
                                    value={password}
                                    onChange={(event) => {
                                        setPassword(event.target.value)
                                    }}
                                    onKeyPress={onKeyPress}
                                />
                                <Button
                                    fullWidth
                                    variant="contained"
                                    sx={{mt: 3, mb: 2}}
                                    onClick={handleSubmit}
                                >
                                    Login
                                </Button>
                            </Box>
                        </Box>
                    </Grid>
                    <Backdrop style={{'zIndex': '100000'}} open={backdropOpen}>
                        <CircularProgress color="inherit"/>
                    </Backdrop>
                </Grid>
            )
        }
    };

    return html();
}

export default withRouter(App);
