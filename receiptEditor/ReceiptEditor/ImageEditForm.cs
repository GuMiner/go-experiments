using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Drawing.Imaging;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace ReceiptEditor
{
    public partial class ImageEditForm : Form
    {
        private Action<ImageAttributes> editCallback;

        public ImageEditForm(Action<ImageAttributes> editCallback)
        {
            InitializeComponent();
            this.editCallback = editCallback;
        }

        private void brightnessTrackbar_Scroll(object sender, EventArgs e)
        {
            ImageAttributes imageAttributes = new ImageAttributes();
            imageAttributes.SetColorMatrix(CreateColorMatrix(brightnessTrackbar.Value), ColorMatrixFlag.Default, ColorAdjustType.Bitmap);

            editCallback(imageAttributes);
        }

        private ColorMatrix CreateColorMatrix(int value)
        {
            float brightness = (float)(value - 50) / 50.0f;
            return new ColorMatrix(
                new[]{
                    new float[] {1,0,0,0,0},
                    new float[] {0,1,0,0,0},
                    new float[] {0,0,1,0,0},
                    new float[] {0,0,0,1,0},
                    new float[] {brightness, brightness, brightness, 0,1}});
        }
    }
}
