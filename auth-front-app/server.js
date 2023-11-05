const express = require('express');
const cookieParser = require('cookie-parser');
const path = require('path');
const app = express();
const port = process.env.PORT || 3000;

app.use(express.json());
app.use(express.urlencoded({extended: false}));
app.use(express.static(path.join(__dirname, 'build')));
app.use(cookieParser());

app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'build', 'index.html'));
});

// TODO: Use a JWT token and ensure it's both signed and encrypted for enhanced security
app.get('/loginRedirect', (req, res) => {
    const accessToken = req.query.token;
    const session = req.query.session;
    res.cookie('access_token', accessToken);
    res.cookie('session', session);
    res.redirect('/');
});

app.listen(port, () => {
    console.log(`Server listening on port ${port}`);
});
