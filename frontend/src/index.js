import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import {SnackbarProvider} from 'notistack';

import {BrowserRouter, Redirect, Route, Switch} from 'react-router-dom';

const theme = createTheme();

const index =
    <ThemeProvider theme={theme}>
        <SnackbarProvider maxSnack={3}>
            <BrowserRouter>
                <Switch>
                    <Route path={"/"} component={App}/>
                    <Redirect from="*" to="/"/>
                </Switch>
            </BrowserRouter>
        </SnackbarProvider>
    </ThemeProvider>

ReactDOM.render(
    index,
    document.getElementById('root')
);
