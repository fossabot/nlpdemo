import React, {useState} from 'react';
import {Button, Checkbox, Form, Grid, Icon, Label, Message, TextArea} from "semantic-ui-react";
import {API_BASE} from "../constants";

const IndexPage = () => {

  const [detectLang, setDetectLanguage] = useState(true);
  const [text, setText] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const [detectedLanguage, setDetectedLanguage] = useState({});
  const [detectedSentiment, setDetectedSentiment] = useState({});

  const submitForm = () => {
    setLoading(true)

    if (detectLang) {
      detectLanguage();
    }

    fetch(API_BASE + '/sentiment', {
      method: 'POST',
      body: JSON.stringify({
        'text': text
      })
    }).then((d) => {
      if (d.status === 200) {
        return d.json();
      } else {
        setError(d.statusText);
      }
    }).then((d) => {
      setDetectedSentiment(d);
    }).finally(() => setLoading(false));

  };

  const detectLanguage = () => {
    fetch(API_BASE + '/language/detect', {
      method: 'POST',
      body: JSON.stringify({
        'text': text
      })
    }).then((d) => {
      if (d.status === 200) {
        return d.json();
      } else {
        setError(d.statusText);
      }
    }).then((d) => {
      setDetectedLanguage(d);
    })
  };

  return (
    <Grid style={{marginTop: 40}}>
      <Grid.Row centered>
        <Grid.Column width={2}/>

        <Grid.Column width={12} stretched>
          <Form error loading={loading}>

            <h3>Sentimental Analysis Demonstration</h3>
            {error && <Message error content={error}/>}

            <Form.Field>
              <label>Enter text</label>
              <TextArea size='large' placeholder='Enter your text' value={text} onChange={(e, d) => setText(d.value)}/>
            </Form.Field>

            <Form.Field>
              <Checkbox label='Detect Language' checked={detectLang} onChange={(e, d) => setDetectLanguage(d.checked)}/>
            </Form.Field>

            <Button onClick={submitForm} primary type='submit' icon labelPosition='right'>Submit<Icon name='right arrow'/></Button>
          </Form>

          <hr/>
          <br/>

          {detectLang && detectedLanguage.language && <div>
            Detected language is <b>{detectedLanguage.language}</b> <Label color='teal'>{detectedLanguage.confidence}</Label>
          </div>}

          {detectedSentiment.sentiment && <div>
            Detected sentiment is <b>{detectedSentiment.sentiment}</b> <Label color='blue'>{detectedSentiment.confidence}</Label>
          </div>}
        </Grid.Column>
        <Grid.Column width={2}/>
      </Grid.Row>
    </Grid>
  );
};

export default IndexPage;
