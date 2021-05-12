import React from 'react'
import { Dialog, DialogType, IDialogProps, Image, Stack, Text, DefaultButton } from '@fluentui/react';
import { ILoginProps } from './Control.types'
import logo_light from '../assets/img/logo_light.svg'
import microsoft_logo from '../assets/img/microsoft-logo.svg'
import github_logo from '../assets/img/github-logo.svg'

export const Login = React.memo<ILoginProps>(() => {

    var pageUrl = encodeURIComponent(window.location.pathname);

    // dialog props
    const props: IDialogProps = {
        hidden: false,
        modalProps: {
            topOffsetFixed: false,
            isBlocking: false,
        },
        dialogContentProps: {
            type: DialogType.normal,
            title: <>
                <Stack horizontalAlign="center">
                    <Image src={logo_light} />
                </Stack>
            </>,
            subText: undefined,
        },
        styles: {
            main: {
                padding: "20px 0"
            }
        },
    };

    return <Dialog {...props}>
        <Stack horizontalAlign="center" tokens={{ childrenGap: 30}}>
            <Text variant="xLarge">Sign in to Pglet</Text>
            <Text variant="medium" style={{ textAlign: "center" }}>You must sign in to access this page. Please continue with one of the options below:</Text>
            <Stack tokens={{ childrenGap: 20}} horizontalAlign="stretch">
                <DefaultButton href={"/api/oauth/github?redirect_url=" + pageUrl} iconProps={{
                    imageProps: {
                        src: github_logo,
                        width: 16,
                        height: 16
                    }
                }} style={{ padding: "0 50px" }}>Sign in with GitHub</DefaultButton>
                <DefaultButton href={"/api/oauth/azure?redirect_url=" + pageUrl} iconProps={{
                    imageProps: {
                        src: microsoft_logo,
                        width: 16,
                        height: 16
                    }
                }}>Sign in with Microsoft account</DefaultButton>
            </Stack>
        </Stack>
    </Dialog>
})