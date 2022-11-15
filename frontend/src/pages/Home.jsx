import React from 'react';
import ReCAPTCHA from 'react-google-recaptcha';
import axios from 'axios';
import { ToastContainer, toast } from 'react-toastify';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';
import * as bech32 from 'bech32';
import { networks, siteKey } from '../config';
import NetworkContext from '../contexts/NetworkContext';

import Section from '../components/Section';
import Button from '../components/Button';
import Icons from '../components/Icons';

import 'react-toastify/dist/ReactToastify.css';
import style from './Home.module.scss';

const validateWalletAddress = (str) => {
  try {
    const { prefix } = bech32.decode(str);
    if (prefix !== 'paloma') {
      throw new Error('Invalid address');
    }
  } catch {
    return 'Enter valid wallet address';
  }
};

const sendSchema = Yup.object().shape({
  address: Yup.string().required('Required'),
  denom: Yup.string().required('Required'),
});

const DENUMS_TO_TOKEN = {
  uluna: 'Luna',
};

const REQUEST_LIMIT_SECS = 30;

class HomeComponent extends React.Component {
  static contextType = NetworkContext;
  recaptchaRef = React.createRef();

  constructor(props) {
    super(props);
    this.state = {
      sending: false,
      verified: false,
      response: '',
    };
  }

  handleCaptcha = (response) => {
    this.setState({
      response,
      verified: true,
    });
  };

  handleSubmit = (values, { resetForm }) => {
    const network = networks.filter(
      (n) => n.chainId === this.context.network
    )[0];
    // same shape as initial values
    this.setState({
      sending: true,
      verified: false,
    });

    this.recaptchaRef.current.reset();

    setTimeout(() => {
      this.setState({ sending: false });
    }, REQUEST_LIMIT_SECS * 1000);

    axios
      .post(network.faucetUrl, {
        address: values.address,
        denom: 'ugrain',
        response: this.state.response,
      })
      .then((res) => {
        let text = res.data;

        if (text === '') {
          toast.success(
            <div>
              <p>Tokens Sent!</p>
            </div>
          );
        } else {
          toast.error(
            <div>
              <p>{text}</p>
            </div>
          );
        }

        // console.log(res);
        // //const response = res.data.response['tx_response'] || res.data.response;
        //
        //
        //   //const url = `https://finder.terra.money/testnet/tx/${response.txhash}`;
        // toast.success(
        //   <div>
        //     <p>
        //       {text}
        //     </p>
        //   </div>
        // );

        resetForm();
      })
      .catch((err) => {
        let errText = err.message;

        if (err.response) {
          if (err.response.data) {
            errText = err.response.data;
          } else {
            switch (err.response.status) {
              case 400:
                errText = 'Invalid request';
                break;
              case 403:
              case 429:
                errText = 'Too many requests';
                break;
              case 404:
                errText = 'Cannot connect to server';
                break;
              case 500:
              case 502:
              case 503:
                errText = 'Faucet service temporary unavailable';
                break;
              default:
                errText = err.message;
            }
          }
        }

        toast.error(`An error occurred: ${errText}`);
      });
  };

  render() {
    return (
      <section className={style.homeContainer}>
        <Section>
          <h2>Paloma Testnet Faucet</h2>
          <article>
            Hello Pigeons! Use this faucet to get GRAIN tokens for the latest
            Paloma Testnest. Plase don’t abuse this service -the number of
            available tokens is limited.
          </article>
          <div className={style.recaptcha}>
            <ReCAPTCHA
              ref={this.recaptchaRef}
              sitekey={siteKey}
              onChange={this.handleCaptcha}
            />
          </div>
          <Formik
            initialValues={{
              address: '',
              denom: 'uluna',
            }}
            validationSchema={sendSchema}
            onSubmit={this.handleSubmit}
          >
            {({ errors, touched }) => (
              <Form className={style.inputContainer}>
                <div className={style.input}>
                  <Field
                    name="address"
                    placeholder="Testnet address"
                    validate={validateWalletAddress}
                  />
                  {errors.address && touched.address ? (
                    <div className="fieldError">{errors.address}</div>
                  ) : null}
                </div>
                <Field type="hidden" name="denom" value="uluna" />
                <Button
                  disabled={!this.state.verified || this.state.sending}
                  type="submit"
                  color="pink"
                >
                  <img src={Icons.RightArrow} />
                  <span>
                    {this.state.sending
                      ? 'Waiting for next tap'
                      : 'Send me tokens'}
                  </span>
                  <img src={Icons.Grain} />
                </Button>
              </Form>
            )}
          </Formik>
        </Section>
        <Section className={style.joinSection}>
          <h2>Need a testnest wallet address?</h2>
          <article>
            Download Nest, Paloma’s Wallet, in the Google Chrome Web Store. Join
            Paloma’s discord server to connect with the flock.
          </article>
          <div className={style.communityButtons}>
            <Button
              color="pink"
              target="_blank"
              href="https://chrome.google.com/webstore/detail/paloma-nestbeta/cjmmdephaciiailjnoikekdebkcbcfmi?hl=en&authuser=1"
            >
              <img src={Icons.Egg} />
              <span>Download Nest Wallet</span>
            </Button>
            <Button
              color="black"
              target="_blank"
              href="https://discord.gg/tNqkNHvVNc"
            >
              <img src={Icons.Discord} />
              <span>Join the Community</span>
            </Button>
          </div>
        </Section>
        <ToastContainer
          position="top-right"
          autoClose={5000}
          hideProgressBar
          newestOnTop
          closeOnClick
          rtl={false}
          pauseOnVisibilityChange
          pauseOnHover
        />
      </section>
    );
  }
}

export default HomeComponent;
